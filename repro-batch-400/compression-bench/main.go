// compression-bench is a small OTLP HTTP receiver that measures decompression
// performance of snappy, zstd, and gzip on real protobuf payloads.
//
// It listens on /v1/logs, /v1/traces, and /v1/metrics, decompresses whatever
// the sender used, then re-compresses + re-decompresses with all three codecs
// to produce a side-by-side comparison.
//
// Usage:
//
//	go run main.go [-port :4318]
package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zstd"
)

// zstd encoder/decoder are reused across requests — they are goroutine-safe
// when constructed with default options.
var (
	zstdEncoder *zstd.Encoder
	zstdDecoder *zstd.Decoder
)

func init() {
	var err error
	zstdEncoder, err = zstd.NewWriter(nil)
	if err != nil {
		log.Fatalf("creating zstd encoder: %v", err)
	}
	zstdDecoder, err = zstd.NewReader(nil)
	if err != nil {
		log.Fatalf("creating zstd decoder: %v", err)
	}
}

func main() {
	port := flag.String("port", ":4318", "listen address (e.g. :4318)")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/logs", makeHandler("logs"))
	mux.HandleFunc("/v1/traces", makeHandler("traces"))
	mux.HandleFunc("/v1/metrics", makeHandler("metrics"))

	log.Printf("compression-bench listening on %s", *port)
	if err := http.ListenAndServe(*port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func makeHandler(signal string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		compressed, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("reading body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		encoding := r.Header.Get("Content-Encoding")

		// Decompress the original payload, measuring only decompression time.
		decompressStart := time.Now()
		raw, err := decompress(encoding, compressed)
		decompressElapsed := time.Since(decompressStart)
		if err != nil {
			http.Error(w, fmt.Sprintf("decompressing %s body: %v", encoding, err), http.StatusBadRequest)
			return
		}

		rawSize := float64(len(raw)) / (1024 * 1024)

		// Benchmark all three codecs: compress then decompress.
		snappyResult := benchCodec("snappy", raw)
		zstdResult := benchCodec("zstd", raw)
		gzipResult := benchCodec("gzip", raw)

		// Print the summary. If there was no encoding on the wire, the
		// decompressElapsed is essentially zero (just a copy); we still print it
		// for consistency.
		fmt.Printf(
			"[%s] raw: %.2f MB | snappy: %.2f MB (%.1fx, compress: %dms, decompress: %dms) | zstd: %.2f MB (%.1fx, compress: %dms, decompress: %dms) | gzip: %.2f MB (%.1fx, compress: %dms, decompress: %dms) | wire encoding: %s decompress: %dms\n",
			signal,
			rawSize,
			snappyResult.compressedMB, rawSize/snappyResult.compressedMB,
			snappyResult.compressMs, snappyResult.decompressMs,
			zstdResult.compressedMB, rawSize/zstdResult.compressedMB,
			zstdResult.compressMs, zstdResult.decompressMs,
			gzipResult.compressedMB, rawSize/gzipResult.compressedMB,
			gzipResult.compressMs, gzipResult.decompressMs,
			encodingLabel(encoding), decompressElapsed.Milliseconds(),
		)

		// Respond with 200 OK and an empty protobuf response body. The OTel
		// Collector expects a valid ExportXxxServiceResponse; an empty body with
		// 200 is sufficient to satisfy most exporters.
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.WriteHeader(http.StatusOK)
	}
}

type codecResult struct {
	compressedMB float64
	compressMs   int64
	decompressMs int64
}

func benchCodec(codec string, raw []byte) codecResult {
	// Compress.
	compressStart := time.Now()
	compressed, err := compress(codec, raw)
	compressElapsed := time.Since(compressStart)
	if err != nil {
		log.Printf("compress %s error: %v", codec, err)
		return codecResult{}
	}

	// Decompress — we discard the result; we only care about timing.
	decompressStart := time.Now()
	_, err = decompress(codec, compressed)
	decompressElapsed := time.Since(decompressStart)
	if err != nil {
		log.Printf("decompress %s error: %v", codec, err)
		return codecResult{}
	}

	return codecResult{
		compressedMB: float64(len(compressed)) / (1024 * 1024),
		compressMs:   compressElapsed.Milliseconds(),
		decompressMs: decompressElapsed.Milliseconds(),
	}
}

// compress encodes src with the named codec and returns the compressed bytes.
func compress(codec string, src []byte) ([]byte, error) {
	switch codec {
	case "snappy":
		// Block format — matches what the OTel Collector sends when
		// compression: snappy is set on the otlphttp exporter.
		return snappy.Encode(nil, src), nil

	case "zstd":
		return zstdEncoder.EncodeAll(src, nil), nil

	case "gzip":
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		if _, err := w.Write(src); err != nil {
			return nil, fmt.Errorf("gzip write: %w", err)
		}
		if err := w.Close(); err != nil {
			return nil, fmt.Errorf("gzip close: %w", err)
		}
		return buf.Bytes(), nil

	default:
		return nil, fmt.Errorf("unknown codec %q", codec)
	}
}

// decompress decodes src using the named encoding. An empty or "identity"
// encoding returns src unchanged.
func decompress(encoding string, src []byte) ([]byte, error) {
	switch encoding {
	case "snappy":
		out, err := snappy.Decode(nil, src)
		if err != nil {
			return nil, fmt.Errorf("snappy decode: %w", err)
		}
		return out, nil

	case "zstd":
		out, err := zstdDecoder.DecodeAll(src, nil)
		if err != nil {
			return nil, fmt.Errorf("zstd decode: %w", err)
		}
		return out, nil

	case "gzip":
		r, err := gzip.NewReader(bytes.NewReader(src))
		if err != nil {
			return nil, fmt.Errorf("gzip reader: %w", err)
		}
		defer r.Close()
		out, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("gzip read: %w", err)
		}
		return out, nil

	case "", "identity":
		// No compression; return a copy so callers can mutate freely.
		out := make([]byte, len(src))
		copy(out, src)
		return out, nil

	default:
		return nil, fmt.Errorf("unsupported Content-Encoding %q", encoding)
	}
}

func encodingLabel(enc string) string {
	if enc == "" {
		return "none"
	}
	return enc
}
