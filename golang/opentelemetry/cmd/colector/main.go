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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

var serviceName = semconv.ServiceNameKey.String("test-service")

// // Initialize a gRPC connection to be used by both the tracer and meter
// // providers.
// func initConn() (*grpc.ClientConn, error) {
// 	// It connects the OpenTelemetry Collector through local gRPC connection.
// 	// You may replace `localhost:4317` with your endpoint.
// 	conn, err := grpc.NewClient("localhost:4317",
// 		// Note the use of insecure transport here. TLS is recommended in production.
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
// 	}

// 	return conn, err
// }

// Initializes an OTLP exporter, and configures the corresponding meter provider.
func initMeterProvider(ctx context.Context, res *resource.Resource) (func(context.Context) error, error) {
	metricExporter, err := otlpmetrichttp.New(
		ctx,
		// otlpmetrichttp.WithEndpoint("/metrics"),
		otlpmetrichttp.WithEndpointURL("http://localhost:4318/metrics"),
		otlpmetrichttp.WithInsecure(),
	)
	// metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
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

	// conn, err := initConn()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// The service name used to display traces in backends
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

	meter := otel.Meter("test-meter")

	srvMetrics := &http.Server{
		Addr: ":4318",
		// Handler: newHTTPHandler(),
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srvMetrics.Handler = mux

	srvErr := make(chan error, 1)
	go func() {
		log.Printf("Starting metrics server on :4318")
		srvErr <- srvMetrics.ListenAndServe()
	}()

	// Attributes represent additional key-value descriptors that can be bound
	// to a metric observer or recorder.
	commonAttrs := []attribute.KeyValue{
		attribute.String("attrA", "chocolate"),
		attribute.String("attrB", "raspberry"),
		attribute.String("attrC", "vanilla"),
	}

	runCount, err := meter.Int64Counter("run", metric.WithDescription("The number of times the iteration ran"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		runCount.Add(ctx, 1, metric.WithAttributes(commonAttrs...))
		log.Printf("Doing really hard work (%d / 10)\n", i+1)

		<-time.After(time.Second)
	}

	// create a srv in the port 4318 to receive the metrics

	select {
	case err = <-srvErr:
		return
	case <-ctx.Done():
		cancel()
	}
}
