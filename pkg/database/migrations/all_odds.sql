DROP TABLE IF EXISTS all_odds;

CREATE TABLE IF NOT EXISTS all_odds (
    id TEXT,
    market TEXT,
    grouping_key TEXT,
    selection TEXT,
    sportsbook TEXT,
    price REAL,
    url TEXT,
    name TEXT
) STRICT;

CREATE INDEX IF NOT EXISTS all_odds_idx ON all_odds(id, market, grouping_key);
