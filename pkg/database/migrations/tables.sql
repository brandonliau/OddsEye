DROP TABLE IF EXISTS all_fixtures;
DROP TABLE IF EXISTS all_odds;
DROP TABLE IF EXISTS arbitrage;
DROP TABLE IF EXISTS expected_value;

CREATE TABLE IF NOT EXISTS all_fixtures (
    id TEXT,
    start_date TEXT,
    home_team TEXT,
    away_team TEXT,
    sport TEXT,
    league TEXT
) STRICT;

CREATE TABLE IF NOT EXISTS all_odds (
    id TEXT,
    market TEXT,
    selection TEXT,
    sportsbook TEXT,
    price REAL,
    url TEXT,
    grouping_key TEXT
) STRICT;

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

CREATE TABLE IF NOT EXISTS expected_value (
    id TEXT,
    market TEXT,
    selection TEXT,
    grouping_key TEXT,
    price TEXT,
    novig_mult REAL,
    novig_add REAL,
    novig_pow REAL,
    novig_shin REAL,
    novig_wc REAL,
    mult_ev REAL GENERATED ALWAYS AS (((1 / novig_mult) * price) - 1) STORED,
    add_ev REAL GENERATED ALWAYS AS (((1 / novig_add) * price) - 1) STORED,
    pow_ev REAL GENERATED ALWAYS AS (((1 / novig_pow) * price) - 1) STORED,
    shin_ev REAL GENERATED ALWAYS AS (((1 / novig_shin) * price) - 1) STORED,
    wc_ev REAL GENERATED ALWAYS AS (((1 / novig_wc) * price) - 1) STORED
) STRICT;

CREATE INDEX IF NOT EXISTS id_all_fixtures_idx ON all_fixtures(id);
CREATE INDEX IF NOT EXISTS id_market_all_odds_idx ON all_odds(id, market, grouping_key);
CREATE INDEX IF NOT EXISTS id_arbitrage_idx ON arbitrage(id, market, grouping_key);
CREATE INDEX IF NOT EXISTS id_expected_value_idx ON expected_value(id, market, grouping_key);
