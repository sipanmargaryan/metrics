package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

// SetupOTLP enables export when an endpoint is configured and returns its cleanup.
func SetupOTLP(service string) (func(), error) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" && os.Getenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT") == "" {
		return func() {}, nil
	}

	logger, shutdown, err := NewOTLP(context.Background(), os.Stdout, service)
	if err != nil {
		return nil, err
	}
	slog.SetDefault(logger)

	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(ctx); err != nil {
			New(os.Stderr, service).Error("flush exported logs", "error", err)
		}
	}, nil
}

// NewOTLP keeps JSON stdout logging and exports logs asynchronously over OTLP/HTTP.
// The exporter reads the standard OTEL_EXPORTER_OTLP_* environment variables.
// Call shutdown after the final log entry to flush pending records.
func NewOTLP(ctx context.Context, w io.Writer, service string) (*slog.Logger, func(context.Context) error, error) {
	exporter, err := otlploghttp.New(ctx, otlploghttp.WithTimeout(3*time.Second))
	if err != nil {
		return nil, nil, err
	}
	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(resource.NewSchemaless(attribute.String("service.name", service))),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)
	handler := slog.NewMultiHandler(
		slog.NewJSONHandler(w, nil),
		otelslog.NewHandler("metrics/internal/logging", otelslog.WithLoggerProvider(provider)),
	)
	return slog.New(handler).With("service_name", service), provider.Shutdown, nil
}
