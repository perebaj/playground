// fixture-gen produces lightweight parquet fixtures that mimic the magnolia
// telemetry schema. The output layout is the org-first hive partitioning
// the design doc proposes:
//
//	{out}/org_id={org}/tables/{signal}/date={YYYY-MM-DD}/data.parquet
//
// We write the same 5 OTLP signal tables magnolia-sampler defines (traces,
// logs, metrics_gauge, metrics_sum, metrics_histogram) so DuckDB can read
// them with the same view definitions we plan to use in production.
package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/parquet-go/parquet-go"
)

// 5-table schema mirroring magnolia-sampler/internal/duckdb/schema.go.
// Attribute columns are JSON-encoded VARCHARs like the real samples.

type Trace struct {
	Timestamp          int64  `parquet:"Timestamp"`
	TraceId            string `parquet:"TraceId"`
	SpanId             string `parquet:"SpanId"`
	ParentSpanId       string `parquet:"ParentSpanId,optional"`
	TraceState         string `parquet:"TraceState,optional"`
	SpanName           string `parquet:"SpanName"`
	SpanKind           string `parquet:"SpanKind"`
	ServiceName        string `parquet:"ServiceName"`
	Duration           uint64 `parquet:"Duration"`
	StatusCode         string `parquet:"StatusCode"`
	StatusMessage      string `parquet:"StatusMessage,optional"`
	ResourceAttributes string `parquet:"ResourceAttributes,optional"`
	ScopeName          string `parquet:"ScopeName,optional"`
	ScopeVersion       string `parquet:"ScopeVersion,optional"`
	SpanAttributes     string `parquet:"SpanAttributes,optional"`
	Events             string `parquet:"Events,optional"`
	Links              string `parquet:"Links,optional"`
}

type Log struct {
	Timestamp          int64  `parquet:"Timestamp"`
	ObservedTimestamp  int64  `parquet:"ObservedTimestamp"`
	TraceId            string `parquet:"TraceId,optional"`
	SpanId             string `parquet:"SpanId,optional"`
	TraceFlags         uint8  `parquet:"TraceFlags,optional"`
	SeverityText       string `parquet:"SeverityText,optional"`
	SeverityNumber     uint8  `parquet:"SeverityNumber"`
	Body               string `parquet:"Body,optional"`
	EventName          string `parquet:"EventName,optional"`
	ServiceName        string `parquet:"ServiceName"`
	ResourceSchemaUrl  string `parquet:"ResourceSchemaUrl,optional"`
	ResourceAttributes string `parquet:"ResourceAttributes,optional"`
	ScopeName          string `parquet:"ScopeName,optional"`
	ScopeVersion       string `parquet:"ScopeVersion,optional"`
	ScopeSchemaUrl     string `parquet:"ScopeSchemaUrl,optional"`
	ScopeAttributes    string `parquet:"ScopeAttributes,optional"`
	LogAttributes      string `parquet:"LogAttributes,optional"`
}

// Metric columns shared by gauge/sum/histogram. Embedded in the three concrete
// metric structs below.
type metricCommon struct {
	ServiceName           string `parquet:"ServiceName"`
	ResourceSchemaUrl     string `parquet:"ResourceSchemaUrl,optional"`
	ResourceAttributes    string `parquet:"ResourceAttributes,optional"`
	ScopeName             string `parquet:"ScopeName,optional"`
	ScopeVersion          string `parquet:"ScopeVersion,optional"`
	ScopeSchemaUrl        string `parquet:"ScopeSchemaUrl,optional"`
	ScopeAttributes       string `parquet:"ScopeAttributes,optional"`
	ScopeDroppedAttrCount uint32 `parquet:"ScopeDroppedAttrCount,optional"`
	MetricName            string `parquet:"MetricName"`
	MetricDescription     string `parquet:"MetricDescription,optional"`
	MetricUnit            string `parquet:"MetricUnit,optional"`
	Attributes            string `parquet:"Attributes,optional"`
	StartTimeUnix         int64  `parquet:"StartTimeUnix,optional"`
	TimeUnix              int64  `parquet:"TimeUnix"`
	Flags                 uint32 `parquet:"Flags,optional"`
	Exemplars             string `parquet:"Exemplars,optional"`
}

type MetricGauge struct {
	metricCommon
	Value float64 `parquet:"Value"`
}

type MetricSum struct {
	metricCommon
	Value                  float64 `parquet:"Value"`
	AggregationTemporality int32   `parquet:"AggregationTemporality"`
	IsMonotonic            bool    `parquet:"IsMonotonic"`
}

type MetricHistogram struct {
	metricCommon
	Count          uint64  `parquet:"Count"`
	Sum            float64 `parquet:"Sum,optional"`
	Min            float64 `parquet:"Min,optional"`
	Max            float64 `parquet:"Max,optional"`
	BucketCounts   string  `parquet:"BucketCounts,optional"`
	ExplicitBounds string  `parquet:"ExplicitBounds,optional"`
}

// orgProfile describes the shape of synthetic telemetry per fake org so the
// generated data resembles plausible diversity (cloud, k8s presence, etc.).
type orgProfile struct {
	id           string
	resourceJSON string
	services     []string
	scopes       []string
}

var orgs = []orgProfile{
	{
		id:           "org_test_alpha",
		resourceJSON: `{"cloud.provider":"gcp","cloud.platform":"gcp_kubernetes_engine","k8s.cluster.name":"alpha-prod","k8s.namespace.name":"default"}`,
		services:     []string{"checkout", "cart", "frontend"},
		scopes:       []string{"io.opentelemetry.spring-boot", "@opentelemetry/instrumentation-express"},
	},
	{
		id:           "org_test_beta",
		resourceJSON: `{"cloud.provider":"aws","cloud.platform":"aws_ec2"}`,
		services:     []string{"payment-api", "auth-service", "recon-worker"},
		scopes:       []string{"go.opentelemetry.io/otel/sdk/tracer", "github.com/XSAM/otelsql"},
	},
	{
		id:           "org_test_demo",
		resourceJSON: `{"cloud.provider":"gcp","cloud.platform":"gcp_kubernetes_engine","k8s.cluster.name":"demo-cluster","k8s.namespace.name":"opentelemetry-demo","telemetry.sdk.name":"opentelemetry"}`,
		services:     []string{"ad", "flagd", "product-catalog"},
		scopes:       []string{"connectrpc.com/otelconnect", "io.opentelemetry.grpc-1.6"},
	},
}

var dates = []string{"2026-06-22", "2026-06-29"}

// Row counts kept small on purpose — this fixture exists to exercise the
// shape, not the scale. Defaults chosen so the six per-signal parquets
// stay under 1 MB combined and load instantly. Override any of them via
// flags to stress specific dimensions (e.g. -traces 60000 to reproduce
// upstream quack extension bugs that only trigger past ~20k rows).
const (
	defaultTracesPerFile    = 50
	defaultLogsPerFile      = 100
	defaultGaugePerFile     = 200
	defaultSumPerFile       = 200
	defaultHistogramPerFile = 50
)

func main() {
	out := flag.String("out", "./fixtures", "output directory for parquet fixtures")
	tracesPerFile := flag.Int("traces", defaultTracesPerFile, "traces rows per parquet file")
	logsPerFile := flag.Int("logs", defaultLogsPerFile, "logs rows per parquet file")
	gaugePerFile := flag.Int("gauge", defaultGaugePerFile, "metrics_gauge rows per parquet file")
	sumPerFile := flag.Int("sum", defaultSumPerFile, "metrics_sum rows per parquet file")
	histogramPerFile := flag.Int("histogram", defaultHistogramPerFile, "metrics_histogram rows per parquet file")
	flag.Parse()

	for _, org := range orgs {
		for _, date := range dates {
			rng := rand.New(rand.NewSource(seed(org.id, date)))
			ts := dateToUnixNanos(date)

			writeParquet(filepath.Join(*out, "org_id="+org.id, "tables", "traces", "date="+date, "data.parquet"),
				traces(rng, org, ts, *tracesPerFile))
			writeParquet(filepath.Join(*out, "org_id="+org.id, "tables", "logs", "date="+date, "data.parquet"),
				logs(rng, org, ts, *logsPerFile))
			writeParquet(filepath.Join(*out, "org_id="+org.id, "tables", "metrics_gauge", "date="+date, "data.parquet"),
				gauges(rng, org, ts, *gaugePerFile))
			writeParquet(filepath.Join(*out, "org_id="+org.id, "tables", "metrics_sum", "date="+date, "data.parquet"),
				sums(rng, org, ts, *sumPerFile))
			writeParquet(filepath.Join(*out, "org_id="+org.id, "tables", "metrics_histogram", "date="+date, "data.parquet"),
				histograms(rng, org, ts, *histogramPerFile))
		}
	}

	log.Printf("wrote fixtures to %s for %d orgs × %d dates × 5 signals (traces=%d/file)",
		*out, len(orgs), len(dates), *tracesPerFile)
}

func traces(rng *rand.Rand, org orgProfile, ts int64, n int) []Trace {
	out := make([]Trace, n)
	for i := 0; i < n; i++ {
		out[i] = Trace{
			Timestamp:          ts + int64(i)*1_000_000_000,
			TraceId:            randHex(rng, 32),
			SpanId:             randHex(rng, 16),
			ParentSpanId:       randHex(rng, 16),
			SpanName:           pick(rng, []string{"GET /api/v1", "POST /checkout", "kafka.consume", "redis.GET"}),
			SpanKind:           pick(rng, []string{"SPAN_KIND_SERVER", "SPAN_KIND_CLIENT", "SPAN_KIND_INTERNAL"}),
			ServiceName:        pick(rng, org.services),
			Duration:           uint64(rng.Int63n(500_000_000)),
			StatusCode:         pick(rng, []string{"STATUS_CODE_OK", "STATUS_CODE_OK", "STATUS_CODE_ERROR"}),
			ResourceAttributes: org.resourceJSON,
			ScopeName:          pick(rng, org.scopes),
			ScopeVersion:       "1.0.0",
			SpanAttributes:     `{"http.method":"GET","http.status_code":200}`,
		}
	}
	return out
}

func logs(rng *rand.Rand, org orgProfile, ts int64, n int) []Log {
	out := make([]Log, n)
	for i := 0; i < n; i++ {
		sev := pick(rng, []string{"INFO", "INFO", "INFO", "WARN", "ERROR"})
		out[i] = Log{
			Timestamp:          ts + int64(i)*500_000_000,
			ObservedTimestamp:  ts + int64(i)*500_000_000,
			SeverityText:       sev,
			SeverityNumber:     severityNumber(sev),
			Body:               pick(rng, []string{"request handled", "cache miss", "connection refused", "user signed in"}),
			ServiceName:        pick(rng, org.services),
			ResourceAttributes: org.resourceJSON,
			ScopeName:          pick(rng, org.scopes),
			ScopeVersion:       "1.0.0",
			LogAttributes:      `{"request.id":"abc123"}`,
		}
	}
	return out
}

func gauges(rng *rand.Rand, org orgProfile, ts int64, n int) []MetricGauge {
	out := make([]MetricGauge, n)
	for i := 0; i < n; i++ {
		out[i] = MetricGauge{
			metricCommon: metricCommon{
				ServiceName:        pick(rng, org.services),
				ResourceAttributes: org.resourceJSON,
				ScopeName:          pick(rng, org.scopes),
				ScopeVersion:       "1.0.0",
				MetricName:         pick(rng, []string{"system.cpu.utilization", "k8s.pod.memory.usage", "process.runtime.go.goroutines"}),
				MetricUnit:         "1",
				TimeUnix:           ts + int64(i)*250_000_000,
			},
			Value: rng.Float64() * 100,
		}
	}
	return out
}

func sums(rng *rand.Rand, org orgProfile, ts int64, n int) []MetricSum {
	out := make([]MetricSum, n)
	for i := 0; i < n; i++ {
		out[i] = MetricSum{
			metricCommon: metricCommon{
				ServiceName:        pick(rng, org.services),
				ResourceAttributes: org.resourceJSON,
				ScopeName:          pick(rng, org.scopes),
				ScopeVersion:       "1.0.0",
				MetricName:         pick(rng, []string{"http.server.requests", "kafka.messages.consumed", "redis.commands.total"}),
				MetricUnit:         "1",
				TimeUnix:           ts + int64(i)*250_000_000,
			},
			Value:                  float64(rng.Int63n(10_000)),
			AggregationTemporality: 2, // CUMULATIVE
			IsMonotonic:            true,
		}
	}
	return out
}

func histograms(rng *rand.Rand, org orgProfile, ts int64, n int) []MetricHistogram {
	out := make([]MetricHistogram, n)
	for i := 0; i < n; i++ {
		count := uint64(rng.Int63n(1000) + 1)
		out[i] = MetricHistogram{
			metricCommon: metricCommon{
				ServiceName:        pick(rng, org.services),
				ResourceAttributes: org.resourceJSON,
				ScopeName:          pick(rng, org.scopes),
				ScopeVersion:       "1.0.0",
				MetricName:         pick(rng, []string{"http.server.duration", "db.client.connections.usage"}),
				MetricUnit:         "ms",
				TimeUnix:           ts + int64(i)*250_000_000,
			},
			Count:          count,
			Sum:            float64(count) * (rng.Float64()*50 + 10),
			Min:            rng.Float64() * 5,
			Max:            rng.Float64()*500 + 50,
			BucketCounts:   "[10,20,15,5]",
			ExplicitBounds: "[5,10,50,100]",
		}
	}
	return out
}

// writeParquet creates dst's parent dirs, then writes rows using parquet-go's
// generic writer. Each call produces one file (one row group is fine at this
// scale).
func writeParquet[T any](dst string, rows []T) {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		log.Fatalf("mkdir %s: %v", filepath.Dir(dst), err)
	}
	f, err := os.Create(dst)
	if err != nil {
		log.Fatalf("create %s: %v", dst, err)
	}
	defer f.Close()

	w := parquet.NewGenericWriter[T](f, parquet.Compression(&parquet.Zstd))
	if _, err := w.Write(rows); err != nil {
		log.Fatalf("write %s: %v", dst, err)
	}
	if err := w.Close(); err != nil {
		log.Fatalf("close %s: %v", dst, err)
	}
}

func dateToUnixNanos(d string) int64 {
	t, err := time.Parse("2006-01-02", d)
	if err != nil {
		log.Fatalf("date %s: %v", d, err)
	}
	return t.UnixNano()
}

// seed derives a deterministic per-(org, date) seed so re-runs produce stable
// fixtures (helpful when diffing output across changes).
func seed(orgID, date string) int64 {
	var s int64
	for _, c := range orgID + date {
		s = s*31 + int64(c)
	}
	return s
}

func pick[T any](rng *rand.Rand, xs []T) T { return xs[rng.Intn(len(xs))] }

func randHex(rng *rand.Rand, n int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, n)
	for i := range b {
		b[i] = hex[rng.Intn(16)]
	}
	return string(b)
}

func severityNumber(text string) uint8 {
	switch text {
	case "TRACE":
		return 1
	case "DEBUG":
		return 5
	case "INFO":
		return 9
	case "WARN":
		return 13
	case "ERROR":
		return 17
	case "FATAL":
		return 21
	}
	return 0
}

var _ = fmt.Sprintf
