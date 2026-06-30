// olive-sim is a small HTTP gateway that mimics what Olive will do in
// production: accept typed REST requests, translate them into SQL with the
// caller's org_id baked into the WHERE clause, send the SQL to the Quack
// server over RPC, and stream JSON results back. Free-form SQL endpoints
// are intentionally absent — this is the public contract.
//
// Per-org enforcement is the entire point: the org_id always comes from
// the URL path and is injected by the server. The caller never sets it
// in a body or query string.
//
// Implementation note: this POC shells out to the system `duckdb` binary
// to act as the Quack client. The first attempt embedded DuckDB via
// go-duckdb v2, but that bundle pins DuckDB 1.4.1 while the Quack
// `quack_query` table function requires DuckDB 1.5+. The gap is captured
// in M1-FINDINGS.md; production Olive would either need a newer DuckDB
// binding or stay on the subprocess model.
package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// orgIDRe enforces the input shape we expect (alphanumeric + underscore).
// A caller cannot smuggle SQL into the org_id segment of the URL because
// we reject anything that doesn't match.
var orgIDRe = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

// attrKeyRe protects the attribute-distribution endpoint from accepting
// anything that isn't an OTel-style key.
var attrKeyRe = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func main() {
	quackURI := envOr("QUACK_URI", "quack:localhost:9494")
	quackToken := envOr("QUACK_TOKEN", "spike_token")
	listen := envOr("LISTEN", ":8080")

	if _, err := exec.LookPath("duckdb"); err != nil {
		log.Fatalf("duckdb binary not in PATH: %v", err)
	}

	srv := &server{quackURI: quackURI, quackToken: quackToken}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", srv.health)
	mux.HandleFunc("GET /api/v1/schema", srv.schema)
	mux.HandleFunc("GET /api/v1/{org_id}/services", srv.services)
	mux.HandleFunc("GET /api/v1/{org_id}/traces/count", srv.tracesCount)
	mux.HandleFunc("GET /api/v1/{org_id}/attributes/{key}/distribution", srv.attrDistribution)

	log.Printf("olive-sim listening on %s, forwarding to %s", listen, quackURI)
	log.Fatal(http.ListenAndServe(listen, mux))
}

type server struct {
	quackURI   string
	quackToken string
}

// runQuack composes a duckdb CLI invocation that loads the Quack extension,
// runs the user's SQL through quack_query, and emits CSV. CSV is the easiest
// stable output format to parse from a subprocess; JSON output would require
// additional extension setup.
func (s *server) runQuack(ctx context.Context, sqlText string) ([]string, [][]string, error) {
	escaped := strings.ReplaceAll(sqlText, "'", "''")
	tokenEsc := strings.ReplaceAll(s.quackToken, "'", "''")
	uriEsc := strings.ReplaceAll(s.quackURI, "'", "''")

	script := fmt.Sprintf(
		`.mode csv
.headers on
LOAD quack;
SELECT * FROM quack_query('%s', '%s', token := '%s');`,
		uriEsc, escaped, tokenEsc,
	)

	cmd := exec.CommandContext(ctx, "duckdb")
	cmd.Stdin = strings.NewReader(script)
	out, err := cmd.Output()
	if err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return nil, nil, fmt.Errorf("duckdb exit %d: %s", ee.ExitCode(), string(ee.Stderr))
		}
		return nil, nil, fmt.Errorf("duckdb: %w", err)
	}

	cols, rows, err := parseCSVWithSuccess(out)
	if err != nil {
		return nil, nil, fmt.Errorf("parse csv: %w (raw: %q)", err, truncate(string(out), 200))
	}
	return cols, rows, nil
}

// parseCSVWithSuccess strips Quack's leading "Success\ntrue" block from the
// output, returns the inner header and rows.
func parseCSVWithSuccess(out []byte) ([]string, [][]string, error) {
	r := csv.NewReader(strings.NewReader(string(out)))
	r.FieldsPerRecord = -1
	var blocks [][][]string
	var cur [][]string
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if len(rec) == 1 && rec[0] == "Success" {
			if len(cur) > 0 {
				blocks = append(blocks, cur)
				cur = nil
			}
		}
		cur = append(cur, rec)
	}
	if len(cur) > 0 {
		blocks = append(blocks, cur)
	}
	if len(blocks) == 0 {
		return nil, nil, errors.New("no rows in duckdb output")
	}
	// The last block is the actual query result.
	last := blocks[len(blocks)-1]
	if len(last) == 0 {
		return nil, nil, errors.New("empty result block")
	}
	cols := last[0]
	var rows [][]string
	if len(last) > 1 {
		rows = last[1:]
	}
	return cols, rows, nil
}

func (s *server) health(w http.ResponseWriter, r *http.Request) {
	if _, _, err := s.runQuack(r.Context(), "SELECT 1"); err != nil {
		http.Error(w, "unhealthy: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) schema(w http.ResponseWriter, r *http.Request) {
	cols, rows, err := s.runQuack(r.Context(), "SHOW TABLES")
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"columns": cols, "rows": rows})
}

// services returns service names + span counts for the org named in the path.
// The org_id is validated and inlined into the WHERE clause; the caller cannot
// influence it via body or query string. This is the typed-endpoint pattern
// the design doc commits to.
func (s *server) services(w http.ResponseWriter, r *http.Request) {
	org, err := validatedOrg(r)
	if err != nil {
		writeError(w, err)
		return
	}
	sqlText := fmt.Sprintf(`
		SELECT ServiceName, count(*) AS spans
		FROM traces
		WHERE org_id = '%s'
		GROUP BY 1 ORDER BY 2 DESC
	`, org)
	respond(w, r, s, sqlText)
}

func (s *server) tracesCount(w http.ResponseWriter, r *http.Request) {
	org, err := validatedOrg(r)
	if err != nil {
		writeError(w, err)
		return
	}
	sinceDays := 30
	if v := r.URL.Query().Get("since_days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			sinceDays = n
		}
	}
	sqlText := fmt.Sprintf(`
		SELECT count(*) AS spans
		FROM traces
		WHERE org_id = '%s' AND date >= now() - INTERVAL %d DAY
	`, org, sinceDays)
	respond(w, r, s, sqlText)
}

func (s *server) attrDistribution(w http.ResponseWriter, r *http.Request) {
	org, err := validatedOrg(r)
	if err != nil {
		writeError(w, err)
		return
	}
	key := r.PathValue("key")
	if !attrKeyRe.MatchString(key) {
		writeError(w, errBadInput("invalid attribute key"))
		return
	}
	sqlText := fmt.Sprintf(`
		SELECT
		  json_extract_string(ResourceAttributes, '$."%s"') AS value,
		  count(*) AS spans
		FROM traces
		WHERE org_id = '%s'
		GROUP BY 1 ORDER BY 2 DESC
	`, key, org)
	respond(w, r, s, sqlText)
}

func validatedOrg(r *http.Request) (string, error) {
	org := r.PathValue("org_id")
	if !orgIDRe.MatchString(org) {
		return "", errBadInput("invalid org_id")
	}
	return org, nil
}

func respond(w http.ResponseWriter, r *http.Request, s *server, sqlText string) {
	cols, rows, err := s.runQuack(r.Context(), sqlText)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"columns": cols, "rows": rows})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

type badInputError struct{ msg string }

func (e *badInputError) Error() string      { return e.msg }
func errBadInput(msg string) *badInputError { return &badInputError{msg} }

func writeError(w http.ResponseWriter, err error) {
	var b *badInputError
	if errors.As(err, &b) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": b.Error()})
		return
	}
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
