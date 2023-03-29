module github.com/vanus-labs/vanus-test

go 1.20

require (
	github.com/cloudevents/sdk-go/v2 v2.13.0
	github.com/fatih/color v1.14.1
	github.com/rs/xid v1.4.0
	github.com/rs/zerolog v1.29.0
	github.com/spf13/cobra v1.6.1
	github.com/thediveo/enumflag/v2 v2.0.2
	github.com/vanus-labs/sdk/golang v0.4.3
	github.com/vanus-labs/vanus/client v0.7.2
	github.com/vanus-labs/vanus/observability v0.7.2
	github.com/vanus-labs/vanus/proto v0.7.2
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
)

require (
	cloud.google.com/go/compute v1.14.0 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/scylladb/go-set v1.0.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/vanus-labs/vanus/pkg v0.7.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.36.4 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.11.2 // indirect
	go.opentelemetry.io/otel/sdk v1.11.2 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/exp v0.0.0-20230213192124-5e25df0256eb // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221202195650-67e5cbc046fd // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace (
	github.com/vanus-labs/sdk/golang => ../sdk/golang
	github.com/vanus-labs/vanus => ../vanus
	github.com/vanus-labs/vanus/client => ../vanus/client
	github.com/vanus-labs/vanus/observability => ../vanus/observability
	github.com/vanus-labs/vanus/pkg => ../vanus/pkg
	github.com/vanus-labs/vanus/proto => ../vanus/proto
	github.com/vanus-labs/vanus/raft => ../vanus/raft
)
