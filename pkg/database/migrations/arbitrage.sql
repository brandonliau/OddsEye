DROP TABLE IF EXISTS arbitrage;

CREATE TABLE IF NOT EXISTS arbitrage (
    id TEXT,
    market TEXT,
    grouping_key TEXT,
    selection_α TEXT,
    selection_β TEXT,
    selection_γ TEXT,
    price_α REAL,
    price_β REAL,
    price_γ REAL,
    total_implied_probability REAL,
    vig REAL
) STRICT;

CREATE INDEX IF NOT EXISTS arbitrage_idx ON arbitrage(id, market, grouping_key);
