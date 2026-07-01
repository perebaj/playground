# M2: QuackScan client-side crash on large server-side aggregations

## TL;DR

Any client query that makes the Quack server aggregate over more than
~20-50k rows of the `traces` view crashes the **client** with an
`InternalException` inside `QuackScan`. The server stays healthy — its
memory doesn't spike, `quack_serve` keeps listening. Streaming small
row sets and aggregations over smaller tables both work fine.

Reproducible locally with a single script:

```
scripts/repro-quackscan-bug.sh
```

Full write-up below. The output belongs verbatim in an upstream issue
against `duckdb/duckdb-quack` if we decide to file one.

## What we found

Rolling out the query engine to our internal Kubernetes cluster with
real magnolia-sampler output (125k rows across 5 parquet files),
`SELECT count(*) FROM q.traces` crashed the client immediately. The
same query executed directly on the pod's DuckDB instance (bypassing
Quack) returned `125000` and change without issue.

## Minimal reproduction

Uses the POC's fixture generator (which now supports a `-traces` flag)
plus the existing `quack-duckdb` Docker image.

```sh
# 1. Generate fixtures with 60k rows per traces file (6 files → 360k rows)
go build -o bin/fixture-gen ./cmd/fixture-gen/
./bin/fixture-gen -traces 60000 -out ./fixtures

# 2. Start quack server in Docker
cat > /tmp/init.sql << 'SQL'
INSTALL quack; INSTALL httpfs;
LOAD quack; LOAD httpfs;
CREATE OR REPLACE VIEW traces AS
  SELECT * FROM read_parquet('/fixtures/org_id=*/tables/traces/date=*/data.parquet', hive_partitioning=true);
CALL quack_serve('quack:[::]:9494', token := 'localtoken', allow_other_hostname := true);
SQL

docker run -d --name quack --platform linux/amd64 \
  -v "$PWD/fixtures:/fixtures:ro" \
  -v /tmp/init.sql:/init/init.sql:ro \
  -p 9494:9494 \
  --entrypoint /bin/sh \
  ghcr.io/ollygarden/magnolia/quack-duckdb:0.0.48 \
  -c 'set -eu; rm -f /tmp/f; mkfifo /tmp/f; tail -f /dev/null > /tmp/f & HOME=/tmp exec /usr/local/bin/duckdb -init /init/init.sql < /tmp/f'

# 3. Client: crash query
duckdb -c "INSTALL quack; LOAD quack;
ATTACH 'localhost:9494' AS q (TYPE quack, TOKEN 'localtoken');
SELECT count(*) FROM (SELECT * FROM q.traces LIMIT 50000);"
```

Result: `InternalException` stack trace ending in `QuackScan`. Changing
`LIMIT 50000` to `LIMIT 20000` returns `20000` cleanly.

## Observed threshold

Swept the LIMIT to bracket the trigger, both locally (this POC) and
against a real S3+httpfs server on Kubernetes. Same threshold both
places:

| Rows LIMIT (server materializes) | Result   |
| -------------------------------- | -------- |
| 100                              | ✅ OK     |
| 10 000                           | ✅ OK     |
| 20 000                           | ✅ OK     |
| 50 000                           | ❌ CRASH  |
| 100 000                          | ❌ CRASH  |

Not binary-searched to a precise cutoff, but the fault line is between
20k and 50k rows on both DuckDB 1.5.4 (server) and DuckDB 1.5.3–1.5.4
(client).

## Stack trace

```
Stack Trace:

0        duckdb::Exception::ToJSON(...)
1        duckdb::InternalException::InternalException(...)
2        duckdb::InternalException::InternalException<
             duckdb::LogicalType const&,
             duckdb::LogicalType const&>(...)
3        duckdb::Vector::Reference(duckdb::Vector const&)
4        duckdb::DataChunk::Reference(duckdb::DataChunk&)
5        duckdb::QuackScan(
             duckdb::ClientContext&,
             duckdb::TableFunctionInput&,
             duckdb::DataChunk&)
6        duckdb::PhysicalTableScan::GetDataInternal(...)
7        duckdb::PipelineExecutor::FetchFromSource(...)
8        duckdb::PipelineExecutor::Execute(...)
9        duckdb::PipelineTask::ExecuteTask(...)
10       duckdb::ExecutorTask::Execute(...)
11       duckdb::TaskScheduler::ExecuteForever(...)
```

The `InternalException` constructor variant with two `LogicalType&`
arguments strongly suggests a type mismatch check firing inside
`Vector::Reference`. The immediate caller is `QuackScan` feeding
`DataChunk::Reference`.

## What we ruled out

Spent an afternoon narrowing this down. Confirmed **not** the cause:

1. **Column type oddities**. `traces` has `Duration UBIGINT NOT NULL`
   and multiple JSON-in-VARCHAR columns (`Events`, `Links`,
   `SpanAttributes`, `ResourceAttributes`). Streaming each individually
   via `SELECT <col> FROM q.traces LIMIT 5` works fine.
   `count(*) FROM q.metrics_histogram` (which also has a `UBIGINT`
   column, `Count`) works. So it isn't a per-column serialization
   issue.

2. **httpfs, hive_partitioning, S3 specifically**. Aggregations against
   the smaller views (`q.metrics_gauge`, `q.metrics_histogram`,
   `q.logs`) — same httpfs/hive machinery, different row counts —
   return correct results. And the local POC uses plain filesystem
   parquet (no httpfs); same crash.

3. **Server-side OOM or crash**. During the failing query, server pod
   memory hovers around 1.0-1.15 GB (limit 2 Gi), never approaches the
   ceiling. Server keeps running, no restarts, logs show `quack_serve`
   still active after the client crash. So the server is fine.

4. **Wire-format skew**. Reproduced with matched client + server (both
   DuckDB 1.5.4, same image). Also reproduced with client 1.5.3 against
   server 1.5.4 — same failure. Not a version mismatch.

5. **Aggregation itself**. `count(*) FROM q.traces WHERE 1=0` returns 0
   cleanly (planner elides the scan). The failure only manifests when
   the server actually reads chunks.

## What remains as the trigger

The only differential we couldn't eliminate is the **number of rows
the server materializes during the scan**. Anything ≤ ~20k rows: OK.
Anything ≥ 50k rows: crash. Behavior is identical across `count(*)`,
`SUM(1)`, `count(<col>)`, and windowed `count(*) OVER () … LIMIT 1`.

Reads like an off-by-something or an unhandled edge in the streaming
RPC that only fires past a threshold — probably a chunk boundary,
buffer boundary, or a type descriptor that only gets renegotiated once
the first N chunks are streamed.

## Environment

- Server image: `ghcr.io/ollygarden/magnolia/quack-duckdb:0.0.48`
  (Debian bookworm-slim + DuckDB 1.5.4 static + `INSTALL quack`,
  `INSTALL httpfs` baked at build time).
- Client: DuckDB 1.5.3 macOS arm64 and DuckDB 1.5.4 linux amd64 (same
  failure on both).
- Quack extension: whatever `INSTALL quack` pulls for DuckDB 1.5.4.
- Data: 6 parquet files, 60k rows each, OTLP `traces` schema
  (VARCHAR + JSON-as-VARCHAR + UBIGINT + BIGINT columns). Same crash
  reproduces with a single 60k-row file — number of files isn't
  material.

Happy to run more targeted repros (specific chunk counts, specific
vector widths) if that helps narrow the fault line.
