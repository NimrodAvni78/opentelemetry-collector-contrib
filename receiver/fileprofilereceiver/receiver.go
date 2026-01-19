// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fileprofilereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fileprofilereceiver"

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/extension/xextension/storage"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/matcher"
)

const (
	storageKey = "processed_files"
)

type fileProfileReceiver struct {
	config        *Config
	settings      receiver.Settings
	consumer      xconsumer.Profiles
	unmarshaler   pprofile.Unmarshaler
	fileMatcher   *matcher.Matcher
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	storageClient storage.Client

	// processedFiles tracks files that have been processed to avoid reprocessing
	processedFiles   map[string]fileInfo
	processedFilesMu sync.RWMutex
}

// fileInfo stores information about a processed file
type fileInfo struct {
	ModTime time.Time `json:"mod_time"`
	Size    int64     `json:"size"`
}

func newFileProfileReceiver(
	settings receiver.Settings,
	cfg *Config,
	consumer xconsumer.Profiles,
) (*fileProfileReceiver, error) {
	criteria := matcher.Criteria{
		Include: cfg.Include,
		Exclude: cfg.Exclude,
	}

	fileMatcher, err := matcher.New(criteria)
	if err != nil {
		return nil, err
	}

	return &fileProfileReceiver{
		config:         cfg,
		settings:       settings,
		consumer:       consumer,
		unmarshaler:    &pprofile.ProtoUnmarshaler{},
		fileMatcher:    fileMatcher,
		processedFiles: make(map[string]fileInfo),
	}, nil
}

func (r *fileProfileReceiver) Start(ctx context.Context, host component.Host) error {
	// Initialize storage client if configured
	if r.config.StorageID != nil {
		ext, found := host.GetExtensions()[*r.config.StorageID]
		if !found {
			return errors.New("storage extension not found")
		}
		storageExt, ok := ext.(storage.Extension)
		if !ok {
			return errors.New("extension is not a storage extension")
		}
		client, err := storageExt.GetClient(ctx, component.KindReceiver, r.settings.ID, "")
		if err != nil {
			return err
		}
		r.storageClient = client

		// Load previously processed files from storage
		if err := r.loadProcessedFiles(ctx); err != nil {
			r.settings.Logger.Warn("Failed to load processed files from storage", zap.Error(err))
		}
	}

	// Create archive directory if needed
	if r.config.CompletionAction == CompletionActionMove {
		if err := os.MkdirAll(r.config.MoveDirectory, 0750); err != nil {
			return err
		}
	}

	ctx, r.cancel = context.WithCancel(ctx)
	r.wg.Add(1)
	go r.pollLoop(ctx)

	return nil
}

func (r *fileProfileReceiver) Shutdown(ctx context.Context) error {
	if r.cancel != nil {
		r.cancel()
	}
	r.wg.Wait()

	// Save processed files to storage
	if r.storageClient != nil {
		if err := r.saveProcessedFiles(ctx); err != nil {
			r.settings.Logger.Warn("Failed to save processed files to storage", zap.Error(err))
		}
		return r.storageClient.Close(ctx)
	}

	return nil
}

func (r *fileProfileReceiver) pollLoop(ctx context.Context) {
	defer r.wg.Done()

	// Process existing files on first run if start_at is "beginning"
	if r.config.StartAt == "beginning" {
		r.poll(ctx)
	}

	ticker := time.NewTicker(r.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.poll(ctx)
		}
	}
}

func (r *fileProfileReceiver) poll(ctx context.Context) {
	files, err := r.fileMatcher.MatchFiles()
	if err != nil {
		r.settings.Logger.Error("Failed to match files", zap.Error(err))
		return
	}

	for _, file := range files {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if r.shouldSkipFile(file) {
			continue
		}

		if err := r.processFile(ctx, file); err != nil {
			r.settings.Logger.Error("Failed to process file",
				zap.String("file", file),
				zap.Error(err))
			continue
		}

		if err := r.handleCompletedFile(file); err != nil {
			r.settings.Logger.Error("Failed to handle completed file",
				zap.String("file", file),
				zap.Error(err))
		}
	}
}

func (r *fileProfileReceiver) shouldSkipFile(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		return true
	}

	r.processedFilesMu.RLock()
	defer r.processedFilesMu.RUnlock()

	if processed, ok := r.processedFiles[file]; ok {
		// Skip if file hasn't changed
		if processed.ModTime.Equal(info.ModTime()) && processed.Size == info.Size() {
			return true
		}
	}

	return false
}

func (r *fileProfileReceiver) processFile(ctx context.Context, file string) error {
	r.settings.Logger.Debug("Processing file", zap.String("file", file))

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	profiles, err := r.unmarshaler.UnmarshalProfiles(data)
	if err != nil {
		return err
	}

	if profiles.ResourceProfiles().Len() == 0 {
		r.settings.Logger.Debug("Skipping empty profile", zap.String("file", file))
		r.markFileProcessed(file)
		return nil
	}

	if err := r.consumer.ConsumeProfiles(ctx, profiles); err != nil {
		return err
	}

	r.markFileProcessed(file)
	r.settings.Logger.Info("Successfully processed profile",
		zap.String("file", file),
		zap.Int("resource_profiles", profiles.ResourceProfiles().Len()))

	return nil
}

func (r *fileProfileReceiver) markFileProcessed(file string) {
	info, err := os.Stat(file)
	if err != nil {
		return
	}

	r.processedFilesMu.Lock()
	defer r.processedFilesMu.Unlock()

	r.processedFiles[file] = fileInfo{
		ModTime: info.ModTime(),
		Size:    info.Size(),
	}
}

func (r *fileProfileReceiver) handleCompletedFile(file string) error {
	switch r.config.CompletionAction {
	case CompletionActionDelete:
		r.settings.Logger.Debug("Deleting processed file", zap.String("file", file))
		if err := os.Remove(file); err != nil {
			return err
		}
		// Remove from tracking since file no longer exists
		r.processedFilesMu.Lock()
		delete(r.processedFiles, file)
		r.processedFilesMu.Unlock()

	case CompletionActionMove:
		destPath := filepath.Join(r.config.MoveDirectory, filepath.Base(file))
		// Add timestamp to avoid collisions
		if _, err := os.Stat(destPath); err == nil {
			ext := filepath.Ext(file)
			base := filepath.Base(file[:len(file)-len(ext)])
			destPath = filepath.Join(r.config.MoveDirectory,
				base+"_"+time.Now().Format("20060102_150405")+ext)
		}
		r.settings.Logger.Debug("Moving processed file",
			zap.String("from", file),
			zap.String("to", destPath))
		if err := os.Rename(file, destPath); err != nil {
			return err
		}
		// Remove from tracking since file moved
		r.processedFilesMu.Lock()
		delete(r.processedFiles, file)
		r.processedFilesMu.Unlock()

	case CompletionActionNone, "":
		// Keep tracking to avoid reprocessing
	}

	return nil
}

func (r *fileProfileReceiver) loadProcessedFiles(ctx context.Context) error {
	data, err := r.storageClient.Get(ctx, storageKey)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	r.processedFilesMu.Lock()
	defer r.processedFilesMu.Unlock()

	return json.Unmarshal(data, &r.processedFiles)
}

func (r *fileProfileReceiver) saveProcessedFiles(ctx context.Context) error {
	r.processedFilesMu.RLock()
	data, err := json.Marshal(r.processedFiles)
	r.processedFilesMu.RUnlock()

	if err != nil {
		return err
	}

	return r.storageClient.Set(ctx, storageKey, data)
}
