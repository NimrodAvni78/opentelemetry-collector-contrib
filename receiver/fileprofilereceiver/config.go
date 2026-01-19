// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fileprofilereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fileprofilereceiver"

import (
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
)

// CompletionAction defines what to do with a file after it has been processed.
type CompletionAction string

const (
	// CompletionActionDelete deletes the file after processing.
	CompletionActionDelete CompletionAction = "delete"
	// CompletionActionMove moves the file to an archive directory after processing.
	CompletionActionMove CompletionAction = "move"
	// CompletionActionNone leaves the file in place but tracks it to avoid reprocessing.
	CompletionActionNone CompletionAction = "none"
)

const (
	defaultPollInterval = 1 * time.Second
)

// Config defines the configuration for the file profile receiver.
type Config struct {
	// Include specifies glob patterns for files to include.
	// At least one pattern is required.
	// Example: ["/var/profiles/*.otlp", "/tmp/profiles/**/*.otlp"]
	Include []string `mapstructure:"include"`

	// Exclude specifies glob patterns for files to exclude.
	// Optional.
	Exclude []string `mapstructure:"exclude"`

	// PollInterval is how often to scan for new files.
	// Default: 1s
	PollInterval time.Duration `mapstructure:"poll_interval"`

	// StartAt specifies where to start reading when the receiver starts.
	// "beginning" processes all existing files, "end" only processes new files.
	// Default: "end"
	StartAt string `mapstructure:"start_at"`

	// CompletionAction specifies what to do with files after processing.
	// Options: "delete", "move", "none"
	// Default: "none"
	CompletionAction CompletionAction `mapstructure:"completion_action"`

	// MoveDirectory is the directory to move files to when completion_action is "move".
	// Required when completion_action is "move".
	MoveDirectory string `mapstructure:"move_directory"`

	// StorageID is the ID of a storage extension to use for tracking processed files.
	// When not specified, tracking is done in memory only (state lost on restart).
	StorageID *component.ID `mapstructure:"storage"`
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if len(c.Include) == 0 {
		return errors.New("'include' must specify at least one glob pattern")
	}

	if c.StartAt != "" && c.StartAt != "beginning" && c.StartAt != "end" {
		return errors.New("'start_at' must be 'beginning' or 'end'")
	}

	switch c.CompletionAction {
	case CompletionActionDelete, CompletionActionMove, CompletionActionNone, "":
		// valid
	default:
		return errors.New("'completion_action' must be 'delete', 'move', or 'none'")
	}

	if c.CompletionAction == CompletionActionMove && c.MoveDirectory == "" {
		return errors.New("'move_directory' is required when completion_action is 'move'")
	}

	if c.PollInterval < 0 {
		return errors.New("'poll_interval' must be non-negative")
	}

	return nil
}

func createDefaultConfig() component.Config {
	return &Config{
		PollInterval:     defaultPollInterval,
		StartAt:          "end",
		CompletionAction: CompletionActionNone,
	}
}
