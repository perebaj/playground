package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/perebaj/playground/golang/httpmetrics/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector(), collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	metrics := metric.NewMetrics(registry)

	mux := http.NewServeMux()

	mux.HandleFunc("/foo", metrics.WrapHandler("/foo", fooHandler))
	mux.HandleFunc("/bar", metrics.WrapHandler("/bar", barHandler))
	mux.HandleFunc("/hello", metrics.WrapHandler("/hello", helloHandler))

	go func() {
		log.Println("Starting metrics server on :8081")
		http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	log.Println("Starting server on", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Foo"))
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Bar"))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
