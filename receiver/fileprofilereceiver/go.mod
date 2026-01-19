module github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fileprofilereceiver

go 1.24.0

require (
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.143.0
	go.opentelemetry.io/collector/component v1.49.1-0.20260115162016-5e41fb551263
	go.opentelemetry.io/collector/consumer/xconsumer v0.143.1-0.20260115162016-5e41fb551263
	go.opentelemetry.io/collector/extension/xextension v0.143.1-0.20260115162016-5e41fb551263
	go.opentelemetry.io/collector/pdata/pprofile v0.143.1-0.20260115162016-5e41fb551263
	go.opentelemetry.io/collector/receiver v1.49.1-0.20260115162016-5e41fb551263
	go.opentelemetry.io/collector/receiver/xreceiver v0.143.1-0.20260115162016-5e41fb551263
	go.opentelemetry.io/otel/metric v1.39.0
	go.opentelemetry.io/otel/trace v1.39.0
	go.uber.org/zap v1.27.1
)

require (
	github.com/bmatcuk/doublestar/v4 v4.9.2 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/elastic/lunes v0.2.0 // indirect
	github.com/hashicorp/go-version v1.8.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.143.0 // indirect
	go.opentelemetry.io/collector/consumer v1.49.1-0.20260115162016-5e41fb551263 // indirect
	go.opentelemetry.io/collector/extension v1.49.1-0.20260115162016-5e41fb551263 // indirect
	go.opentelemetry.io/collector/featuregate v1.49.1-0.20260115162016-5e41fb551263 // indirect
	go.opentelemetry.io/collector/internal/componentalias v0.0.0-20260115162016-5e41fb551263 // indirect
	go.opentelemetry.io/collector/pdata v1.49.1-0.20260115162016-5e41fb551263 // indirect
	go.opentelemetry.io/collector/pipeline v1.49.1-0.20260115162016-5e41fb551263 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage => ../../extension/storage

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../internal/coreinternal

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza => ../../pkg/stanza

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest => ../../pkg/pdatatest

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil => ../../pkg/pdatautil

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden => ../../pkg/golden

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common

// Can be removed after 0.144.0 release
replace go.opentelemetry.io/collector/internal/componentalias => go.opentelemetry.io/collector/internal/componentalias v0.0.0-20260115162016-5e41fb551263
