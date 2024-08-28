// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Example using OTLP exporters + collector + third-party backends. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

var serviceName = semconv.ServiceNameKey.String("test-service")

// Initializes an OTLP exporter, and configures the corresponding meter provider.
func initMeterProvider(ctx context.Context, res *resource.Resource) (func(context.Context) error, error) {
	metricExporter, err := otlpmetrichttp.New(
		ctx,
		// otlpmetrichttp.WithEndpoint("/metrics"),
		otlpmetrichttp.WithEndpointURL("http://localhost:7000/metrics"),
		// insecure means that the TLS verification is skipped.
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	return meterProvider.Shutdown, nil
}

func main() {
	log.Printf("Waiting for connection...")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			serviceName,
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	shutdownMeterProvider, err := initMeterProvider(ctx, res)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdownMeterProvider(ctx); err != nil {
			log.Fatalf("failed to shutdown MeterProvider: %s", err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srvMetrics := &http.Server{
		Addr:    ":7000",
		Handler: mux,
	}

	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Starting metrics server on :7000")
		srvErr <- srvMetrics.ListenAndServe()
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: newHTTPHandler(),
	}

	go func() {
		log.Printf("Starting server on :8080")
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-srvErr:
		log.Fatalf("server error: %v", err)
	case <-ctx.Done():
		log.Printf("Shutting down server...")
	}
}

func tmp() {
	provider := metric.NewMeterProvider()
	meter := provider.Meter("service")

	counter, err := meter.Int64Counter("http_requests_total",
	metric.WithDescription("Total number of HTTP requests"), metric.WithUnit("1"),
	metric.WithInstrumentationName("http"),)

	if err != nil {
		log.Fatalf("failed to create counter: %v", err)
	}

	counter.Add(ctx, 1)

	histogram, err := meter.NewInt64Histogram("request_duration_seconds")
	if err != nil {
		log.Fatalf("failed to create histogram: %v", err)
	}

	histogram.Record(ctx, 1, meter.NewLabelSet())



}


func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	handlerFunc := func(path string, f func(http.ResponseWriter, *http.Request)) {
		metric.NewMeterProvider().Meter("service").NewInt64Counter("http_requests_total"
		).AddCallback(func(ctx context.Context, result metric.Int64ObserverResult) {
			result.Observe(1)
		})

		otelHandler := otelhttp.WithRouteTag(path, http.HandlerFunc(f))
		mux.Handle(path, otelHandler)
	}

	handlerFunc("/hello", helloHandler)
	handlerFunc("/jj", jjIsAwesomeHandler)
	handler := otelhttp.NewHandler(mux, "service")

	return handler
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Hello, world!"))
}

func jjIsAwesomeHandler(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(1 * time.Second)
	_, _ = w.Write([]byte("JJ is awesome!"))
}
