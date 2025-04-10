CREATE TABLE IF NOT EXISTS backtest (
    id TEXT,
    market TEXT,
    grouping_key TEXT,
    selection TEXT,
    price REAL,
    bet_type TEXT,
    name TEXT
) STRICT;

CREATE INDEX IF NOT EXISTS backtest_idx ON backtest(id, market, grouping_key);
