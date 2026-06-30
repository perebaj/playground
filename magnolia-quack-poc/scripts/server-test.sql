LOAD quack;
INSTALL httpfs; LOAD httpfs;
CREATE OR REPLACE VIEW traces AS
  SELECT * FROM read_parquet('fixtures/org_id=*/tables/traces/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW logs AS
  SELECT * FROM read_parquet('fixtures/org_id=*/tables/logs/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW metrics_gauge AS
  SELECT * FROM read_parquet('fixtures/org_id=*/tables/metrics_gauge/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW metrics_sum AS
  SELECT * FROM read_parquet('fixtures/org_id=*/tables/metrics_sum/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW metrics_histogram AS
  SELECT * FROM read_parquet('fixtures/org_id=*/tables/metrics_histogram/date=*/data.parquet', hive_partitioning=true);
CALL quack_serve('quack:localhost:9494', token := 'spike_token');
