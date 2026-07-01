# M2: `QuackScan` client-side crash on column-narrowing projections past a row threshold

## TL;DR

Any client query that asks the Quack server to return a `DataChunk` with **fewer columns than the underlying view** — projection, `count(*)`, aggregation, `SELECT c1 FROM two_col_view` — crashes the **client** with an `InternalException` inside `QuackScan` once the underlying data crosses a row-count threshold. The server stays healthy. Streaming the full-width schema (`SELECT * FROM view`) works at every volume we tested (up to 360k rows).

Minimal repro is **two BIGINT columns and ~26k rows** — no traces schema, no S3, no httpfs, no JSON, no UBIGINT.

Reproducible locally with a single script:

```
scripts/repro-quackscan-bug.sh
```

## Minimal reproduction

```sql
-- 1. Trivial parquet: two BIGINT columns, one row per range value.
COPY (SELECT range::BIGINT AS c1, range::BIGINT AS c2 FROM range(26624))
  TO 'p.parquet' (FORMAT PARQUET);
```

Server (any Quack instance):

```sql
LOAD quack; LOAD httpfs;
CREATE OR REPLACE VIEW v AS SELECT * FROM read_parquet('p.parquet');
CALL quack_serve('quack:[::]:9494', token := 't', allow_other_hostname := true);
```

Client:

```sql
INSTALL quack; LOAD quack;
ATTACH 'localhost:9494' AS q (TYPE quack, TOKEN 't');

-- Works: full-schema passthrough
SELECT * FROM q.v;
SELECT c1, c2 FROM q.v;
SELECT c1 + c2 AS s FROM q.v;

-- Crashes: response DataChunk has fewer columns than the source view
SELECT c1 FROM q.v;        -- InternalException in QuackScan
SELECT c2 FROM q.v;        -- InternalException in QuackScan
SELECT count(*) FROM q.v;  -- InternalException in QuackScan
```

## Observed pass/crash matrix

All against the same `v` view over a 2-BIGINT-column parquet, `n` rows.

| Query                                       | Row count | Response columns | Result   |
| ------------------------------------------- | --------- | ---------------- | -------- |
| `SELECT * FROM q.v`                         | 26 624    | 2 (=source)      | ✅ OK     |
| `SELECT c1, c2 FROM q.v`                    | 26 624    | 2 (=source)      | ✅ OK     |
| `SELECT c1 + c2 AS s FROM q.v`              | 26 624    | 1 (derived)      | ✅ OK     |
| `SELECT c1, c2 FROM q.v WHERE c1 > -1`      | 26 624    | 2 (=source)      | ✅ OK     |
| `SELECT c1 FROM q.v`                        | 26 624    | 1 (< source)     | ❌ CRASH  |
| `SELECT c2 FROM q.v`                        | 26 624    | 1 (< source)     | ❌ CRASH  |
| `SELECT c1 FROM q.v WHERE c1 > 0`           | 26 624    | 1 (< source)     | ❌ CRASH  |
| `SELECT count(*) FROM q.v`                  | 26 624    | 1 (aggregate)    | ❌ CRASH  |
| `SELECT c1 FROM q.v` (small)                | 22 528    | 1 (< source)     | ✅ OK     |
| `SELECT count(*) FROM q.v` (small)          | 22 528    | 1 (aggregate)    | ✅ OK     |

Observations:

- The trigger is **response columns < source view columns**, not aggregation per se. `SELECT c1` and `SELECT count(*)` both fail; `SELECT c1 + c2` and `SELECT *` both pass.
- The trigger is **not** column projection alone. It requires a row-count threshold too: identical projections work fine on the same schema at 22 528 rows.
- Below the threshold, every query on our two-column parquet succeeds.

## Row-count threshold sweep

Two BIGINT columns, varying row counts. DuckDB's `STANDARD_VECTOR_SIZE` is 2048 (that's why we sweep in 2048-row multiples):

| Rows (n × 2048 chunks) | `count(*)` | Notes |
| ---------------------- | ---------- | ----- |
| 20 480 (10 chunks)     | ✅ OK       |       |
| 22 528 (11 chunks)     | ✅ OK       |       |
| 24 576 (12 chunks)     | ✅ OK       |       |
| 26 624 (13 chunks)     | ❌ CRASH    |       |
| 28 672 (14 chunks)     | ❌ CRASH    |       |
| 60 000                 | ❌ CRASH    |       |

Fault line lands between 24 576 (12 chunks) and 26 624 (13 chunks) for two BIGINT columns.

Three BIGINT columns: `count(*) FROM q.v` on a 20 480-row parquet passes; 60 000-row parquet crashes — threshold is higher than 20 480 rows for 3-col, so it's not just chunk count. Whatever the fault line's shape, we didn't binary-search across (rows × columns) combinations exhaustively.

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

The `InternalException` constructor variant carrying two `LogicalType&` arguments strongly suggests a type-mismatch check firing inside `Vector::Reference`. The immediate caller is `QuackScan` handing a `DataChunk` into `Reference`. Best guess: the destination `DataChunk` on the client side was allocated with the source view's schema (2 cols) while the incoming chunk has a projected schema (1 col), and `Reference` refuses to point a v[0] slot of type `BIGINT` at a source that has a different column layout.

## What we ruled out

Spent an afternoon narrowing this down. Confirmed **not** the cause:

1. **Any specific column type.** Bug reproduces on plain `BIGINT`/`BIGINT` — no UBIGINT, no JSON-in-VARCHAR, no nullability, no nested types.
2. **The traces schema specifically.** Same crash on `logs` (which has no UBIGINT), and on our synthesized two-column BIGINT parquet.
3. **httpfs, S3, or hive_partitioning.** Reproduces with a plain local `read_parquet('/path/file.parquet')` in Docker.
4. **Number of parquet files.** Reproduces with a single parquet file.
5. **Row count alone.** `SELECT * FROM q.range_view` (built-in `range()` table function, single BIGINT, no file) at 60 000 rows works. `SELECT * FROM q.v1_bigint` at 60 000 rows works. It's row count *combined with* column narrowing.
6. **Aggregation semantics.** `SELECT c1 + c2 AS s FROM q.v` at 26 624 rows works (single-column output, expression). `SELECT c1 FROM q.v` at the same volume crashes. It's not aggregation, it's schema narrowing.
7. **Server-side OOM.** During the failing query, container memory stays flat at ~150 MB, well below the 512 Mi limit. Server keeps running, no restarts, `quack_serve` still listed as active in the server logs.
8. **Wire-format skew.** Reproduces with matched client + server (both DuckDB 1.5.4, same image). Also reproduces with client 1.5.3 against server 1.5.4.
9. **`WHERE 1=0` short-circuit avoids the crash** — the planner elides the scan entirely, so the failure only manifests when the server actually reads chunks.

## What remains as the trigger

The single differential we couldn't eliminate is a combination:

- **The response `DataChunk` has fewer columns than the source view's schema** — either via explicit projection, filtering with projection, or aggregation.
- **The server has to stream more than roughly 12–13 default-width chunks** to satisfy the query (~24 576–26 624 rows for 2-column data). Fault line shifts with column count in ways we didn't fully characterize.

Both conditions together: crash. Either alone: no crash. Full-schema streaming (`SELECT *`) is unaffected at any tested volume.

## Environment

- Server image: `ghcr.io/ollygarden/magnolia/quack-duckdb:0.0.48` (Debian bookworm-slim + DuckDB 1.5.4 static + `INSTALL quack`, `INSTALL httpfs` baked at build time).
- Client: DuckDB 1.5.3 macOS arm64 and DuckDB 1.5.4 linux amd64 (same failure on both).
- Quack extension: whatever `INSTALL quack` pulls for DuckDB 1.5.4.
- Data: single-file parquet, 2 columns × 26 624 rows of `BIGINT`. Same behavior on multi-file hive-partitioned views.

Happy to run more targeted repros (specific chunk widths, other projected column combinations, alternative `Reference()`-triggering statements) if that helps narrow the fault line further.
