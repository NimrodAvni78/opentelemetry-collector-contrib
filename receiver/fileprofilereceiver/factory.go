// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package fileprofilereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fileprofilereceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/xreceiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fileprofilereceiver/internal/metadata"
)

// NewFactory creates a factory for the file profile receiver.
func NewFactory() receiver.Factory {
	return xreceiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		xreceiver.WithProfiles(createProfilesReceiver, metadata.ProfilesStability),
	)
}

func createProfilesReceiver(
	_ context.Context,
	settings receiver.Settings,
	cfg component.Config,
	consumer xconsumer.Profiles,
) (xreceiver.Profiles, error) {
	rCfg := cfg.(*Config)
	return newFileProfileReceiver(settings, rCfg, consumer)
}
