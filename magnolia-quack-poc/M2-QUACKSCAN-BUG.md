# M2: `QuackScan` client-side crash on projections from VIEWs backed by table functions

## TL;DR

Any client query that projects fewer columns than the source view **AND** the source view wraps a table function like `read_parquet` crashes the **client** with an `InternalException` inside `QuackScan` past a row-count threshold. The identical fixture exposed as a materialized `TABLE` instead of a `VIEW` has no such threshold and returns correct results at every volume we tested (up to 500 000 rows).

- Same data.
- Same 2-column BIGINT schema.
- Same client query.
- `CREATE VIEW v AS SELECT * FROM read_parquet(...)` → **CRASH**.
- `CREATE TABLE t AS SELECT * FROM read_parquet(...)` → **OK**.

Minimal repro is **two BIGINT columns, ~26k rows, exposed as a VIEW**. Local script: `scripts/repro-quackscan-bug.sh`.

## Minimal reproduction

Fixture:

```sql
COPY (SELECT range::BIGINT AS c1, range::BIGINT AS c2 FROM range(26624))
  TO 'p.parquet' (FORMAT PARQUET);
```

Server (any Quack instance):

```sql
LOAD quack;
CREATE OR REPLACE VIEW  v AS SELECT * FROM read_parquet('p.parquet');
CREATE OR REPLACE TABLE t AS SELECT * FROM read_parquet('p.parquet');
CALL quack_serve('quack:[::]:9494', token := 't', allow_other_hostname := true);
```

Client:

```sql
INSTALL quack; LOAD quack;
ATTACH 'localhost:9494' AS q (TYPE quack, TOKEN 't');

-- VIEW-backed: crashes on column narrowing past threshold
SELECT c1        FROM q.v;   -- CRASH
SELECT count(*)  FROM q.v;   -- CRASH

-- TABLE-backed: identical queries succeed
SELECT c1        FROM q.t;   -- OK
SELECT count(*)  FROM q.t;   -- OK

-- VIEW-backed with full schema passthrough: OK
SELECT *          FROM q.v;
SELECT c1, c2     FROM q.v;
SELECT c1 + c2    FROM q.v;
```

## Pass / crash matrix

Same 26 624-row 2-BIGINT parquet, exposed as both `v` (VIEW) and `t` (TABLE).

| Query                                    | Response columns | `v` (VIEW) | `t` (TABLE) |
| ---------------------------------------- | ---------------- | ---------- | ----------- |
| `SELECT *`                               | 2 (source)       | ✅ OK       | ✅ OK        |
| `SELECT c1, c2`                          | 2 (source)       | ✅ OK       | ✅ OK        |
| `SELECT c1 + c2 AS s`                    | 1 (derived)      | ✅ OK       | ✅ OK        |
| `SELECT c1`                              | 1 (< source)     | ❌ CRASH    | ✅ OK        |
| `SELECT c2`                              | 1 (< source)     | ❌ CRASH    | ✅ OK        |
| `SELECT c1 WHERE c1 > 0`                 | 1 (< source)     | ❌ CRASH    | ✅ OK        |
| `SELECT count(*)`                        | 1 (aggregate)    | ❌ CRASH    | ✅ OK        |

TABLE-backed queries scale further too — we tested 500 000 rows across 3 columns (BIGINT + BIGINT + VARCHAR) with `count(*)`, projection, and `GROUP BY`; all succeeded.

## Row-count threshold sweep (2-BIGINT VIEW)

`STANDARD_VECTOR_SIZE = 2048`, so we swept at multiples:

| Rows (n × 2048)     | `count(*) FROM q.view` |
| ------------------- | ---------------------- |
| 20 480 (10 chunks)  | ✅ OK                   |
| 22 528 (11 chunks)  | ✅ OK                   |
| 24 576 (12 chunks)  | ✅ OK                   |
| 26 624 (13 chunks)  | ❌ CRASH                |
| 28 672 (14 chunks)  | ❌ CRASH                |
| 60 000              | ❌ CRASH                |

Fault line between 24 576 and 26 624 rows for two BIGINT columns via VIEW. Three-column BIGINT VIEW: 20 480 rows still passes, 60 000 rows crashes — so the threshold shifts with column count.

Every TABLE-backed test we ran passed regardless of row count.

## Stack trace on the client

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

The `InternalException` constructor variant carrying two `LogicalType&` arguments strongly suggests a type-mismatch check firing inside `Vector::Reference`. The immediate caller is `QuackScan` feeding a `DataChunk::Reference`. Best guess (unverified): views wrapping `read_parquet` and materialized tables take different code paths in `QuackScan`'s response-chunk allocation, and the view path mishandles column narrowing past a chunk-count boundary.

## What we ruled out

Confirmed **not** the cause:

1. **Any specific column type.** Bug reproduces on plain `BIGINT`/`BIGINT`. No UBIGINT, no JSON-in-VARCHAR, no nullability, no nested types.
2. **Aggregation semantics.** `SELECT c1 + c2 AS s FROM q.v` (single-column derived output) works at 26 624 rows. `SELECT c1 FROM q.v` (single-column projection) crashes at the same volume.
3. **Row count alone.** `SELECT * FROM q.v` (full-schema passthrough) works at 60 000 rows. TABLE-backed projections work at 500 000 rows.
4. **httpfs, S3, or hive_partitioning.** Reproduces with plain local `read_parquet('file.parquet')` in Docker.
5. **Number of parquet files.** Reproduces with a single file.
6. **Server-side OOM.** During failing queries, server container memory holds flat at ~150 MB (limit 512 Mi). Server keeps running, no restarts.
7. **Wire-format skew.** Reproduces with matched client + server (both DuckDB 1.5.4). Also reproduces with client 1.5.3 against server 1.5.4.

## What remains as the trigger

Two conditions **both required**:

1. The source in the **server-side** catalog is a **VIEW**. Any view — even one that wraps a materialized TABLE, not just those wrapping `read_parquet`.
2. The client query returns **fewer columns than the source view** — projection, `count(*)`, aggregation. Full-schema passthrough (`SELECT *`) is unaffected.

Above a row-count threshold that scales with column count, both together crash.

**Client-side views are unaffected**. A client that runs `CREATE VIEW my_view AS SELECT * FROM q.some_server_table` and then `SELECT c1 FROM my_view` works correctly at every volume we tested. The bug is specifically about server-side view expansion during the quack RPC scan, not about the view abstraction itself.

## Practical implication for consumers

For anyone deploying Quack as a query gateway over parquet-on-S3 (which is our production use case), the natural pattern is `CREATE VIEW … AS SELECT * FROM read_parquet('s3://…')` so the server re-reads S3 on every query and picks up newly-written data. This pattern is exactly what triggers the bug.

Workarounds until upstream fixes it:

- **Expose materialized TABLEs on the server, not views** (`CREATE TABLE t AS SELECT * FROM read_parquet('s3://…')`). Fully avoids the crash but freezes the data snapshot; requires a refresh mechanism (pod restart, cronjob issuing `CREATE OR REPLACE TABLE`, or an S3 notification triggering a reload). Clients can still wrap those tables in their own **client-side** views for naming/semantics — client-side views are unaffected by the bug.
- **Restrict client queries to full-schema passthrough** (`SELECT * FROM q.view`) and filter/aggregate on the client. Poor UX for a query gateway; also transfers all rows over the RPC.
- **Move to a lower row count per view** — either pre-partition or add filtering — so no single query crosses the threshold. Only viable when the schema allows narrow partitions.

## Environment

- Server image: `ghcr.io/ollygarden/magnolia/quack-duckdb:0.0.48` (Debian bookworm-slim + DuckDB 1.5.4 static + `INSTALL quack`, `INSTALL httpfs` baked at build time).
- Client: DuckDB 1.5.3 macOS arm64 and DuckDB 1.5.4 linux amd64 (same failure on both).
- Quack extension: whatever `INSTALL quack` pulls for DuckDB 1.5.4.
- Data: single-file parquet, 2 columns × 26 624 rows of `BIGINT`. Bug also reproduces on multi-file hive-partitioned views and on our production `traces` schema — narrowing to two-column BIGINT was the minimization step.

Happy to run more targeted repros (specific chunk widths, other projected column combinations, alternative table functions instead of `read_parquet`) if that helps narrow the fault line further.
