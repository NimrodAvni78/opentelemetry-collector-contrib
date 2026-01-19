// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate mdatagen metadata.yaml

// Package fileprofilereceiver scans the filesystem for profiling data in OTLP
// binary protobuf format and forwards it to the next component in the pipeline.
package fileprofilereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fileprofilereceiver"
