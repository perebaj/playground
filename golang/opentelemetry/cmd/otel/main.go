package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func rolldice(w http.ResponseWriter, r *http.Request) {
	roll := 1 + rand.Intn(6)

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
}

// SetupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := NewPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up meter provider.
	exp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		handleErr(fmt.Errorf("failed to create OTLP exporter: %w", err))
		return
	}

	meterProvider, err := NewMeterProvider(exp)
	if err != nil {
		handleErr(fmt.Errorf("failed to create meter provider: %w", err))
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func NewPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func NewMeterProvider(exp *otlpmetricgrpc.Exporter) (*metric.MeterProvider, error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("server"),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exp,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second),
		)),
		metric.WithResource(res),
	)
	return meterProvider, nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := SetupOTelSDK(ctx)

	if err != nil {
		return fmt.Errorf("failed to set up OpenTelemetry SDK: %w", err)
	}

	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
	}()

	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(listener net.Listener) context.Context { return ctx },
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      newHTTPHandler(),
	}

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		stop()
	}

	err = srv.Shutdown(ctx)
	return
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	handleFunc := func(path string, f func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(path, http.HandlerFunc(f))
		mux.Handle(path, handler)
	}

	handleFunc("/rolldice", rolldice)
	handler := otelhttp.NewHandler(mux, "server")
	return handler
}
