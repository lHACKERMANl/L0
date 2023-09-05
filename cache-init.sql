CREATE ROLE cacheUser LOGIN PASSWORD '123';

CREATE TABLE IF NOT EXISTS cache_table (
    key TEXT PRIMARY KEY,
    value TEXT
);

CREATE INDEX IF NOT EXISTS idx_cache_key ON cache_table (key);