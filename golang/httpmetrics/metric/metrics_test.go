package metric

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestWrapHandler(t *testing.T) {
	register := prometheus.NewRegistry()
	metric := NewMetrics(register)

	if metric == nil {
		t.Error("Expected metric to be created")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", metric.WrapHandler("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})))

	mux.HandleFunc("/nometric", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	cnt := testutil.CollectAndCount(metric.httpDuration, "http_request_duration_seconds")
	if cnt != 1 {
		t.Errorf("Expected 1, got %d", cnt)
	}

	cnt2 := testutil.CollectAndCount(metric.httpRequestsTotal, "http_requests_total")
	if cnt2 != 1 {
		t.Errorf("Expected 1, got %d", cnt2)
	}

	req2 := httptest.NewRequest(http.MethodGet, "/nometric", nil)
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)

	cnt3 := testutil.CollectAndCount(metric.httpDuration, "http_request_duration_seconds")
	if cnt3 == 0 {
		t.Errorf("Expected 0, got %d", cnt3)
	}
}

func TestNewMetrics(t *testing.T) {
	register := prometheus.NewRegistry()
	metric := NewMetrics(register)

	metric.httpDuration.WithLabelValues("200", "GET", "/").Observe(0.5)

	cnt := testutil.CollectAndCount(metric.httpDuration, "http_request_duration_seconds")
	if cnt != 1 {
		t.Errorf("Expected 1, got %d", cnt)
	}

	metric.httpRequestsTotal.WithLabelValues("200", "GET", "/").Inc()

	cnt2 := testutil.CollectAndCount(metric.httpRequestsTotal, "http_requests_total")
	if cnt2 != 1 {
		t.Errorf("Expected 1, got %d", cnt2)
	}
}
