# Reuse the official duckdb/duckdb binary (multi-arch, maintained by DuckDB
# Labs) layered onto a minimal debian for the shell + cert bundle we need to
# keep the server alive and reach the extension registry at build time.
#
# Why not duckdb/duckdb directly: that image is a scratch container with only
# the binary as the entrypoint. We need a shell to chain init.sql with a
# tail -f /dev/null so duckdb's main process stays alive after CALL
# quack_serve() returns. We also bake the quack and httpfs extensions at
# build time so the pod has no outbound dependency at boot.

FROM duckdb/duckdb:1.5.4 AS duckdb

FROM debian:bookworm-slim

RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates \
 && rm -rf /var/lib/apt/lists/*

COPY --from=duckdb /duckdb /usr/local/bin/duckdb

RUN duckdb -c "INSTALL quack; INSTALL httpfs;"

ENTRYPOINT ["/usr/local/bin/duckdb"]
