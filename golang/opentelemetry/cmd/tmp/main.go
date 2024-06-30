package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func main() {
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	exp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("jojoapi"),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exp,
			metric.WithInterval(3*time.Second))),
		metric.WithResource(resources),
	)

	defer func() {
		if err := meterProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	otel.SetMeterProvider(meterProvider)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: newHTTPHandler(),
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

	// From here, the meterProvider can be used by instrumentation to collect
	// telemetry.
}

func newHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	handleFunc := func(path string, f func(http.ResponseWriter, *http.Request)) {
		handler := otelhttp.WithRouteTag(path, http.HandlerFunc(f))
		mux.Handle(path, handler)
	}

	handleFunc("/hello", helloHandler)
	handler := otelhttp.NewHandler(mux, "server")
	return handler
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
