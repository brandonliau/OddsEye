DROP TABLE IF EXISTS all_fixtures;
DROP TABLE IF EXISTS all_odds;
DROP TABLE IF EXISTS fair_odds;
DROP TABLE IF EXISTS arbitrage;
DROP TABLE IF EXISTS positive_ev;

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

CREATE TABLE IF NOT EXISTS fair_odds (
    id TEXT,
    market TEXT,
    selection TEXT,
    price TEXT,
    novig_mult REAL,
    novig_add REAL,
    novig_pow REAL,
    novig_shin REAL,
    novig_wc REAL,
    grouping_key TEXT
) STRICT;

CREATE TABLE IF NOT EXISTS arbitrage (
    id TEXT,
    market TEXT,
    selection_α TEXT,
    selection_β TEXT,
    selection_γ TEXT,
    price_α REAL,
    price_β REAL,
    price_γ REAL,
    total_implied_probability REAL,
    vig REAL
) STRICT;

CREATE TABLE IF NOT EXISTS positive_ev (
    id TEXT,
    market TEXT,
    selection TEXT,
    price REAL,
    fair_price REAL,
    ev REAL
) STRICT;

CREATE INDEX IF NOT EXISTS id_all_fixtures_idx ON all_fixtures(id);
CREATE INDEX IF NOT EXISTS id_market_all_odds_idx ON all_odds(id, market);
CREATE INDEX IF NOT EXISTS id_fair_odds_idx ON fair_odds(id, market);
CREATE INDEX IF NOT EXISTS id_arbitrage_idx ON arbitrage(id);
CREATE INDEX IF NOT EXISTS id_arbitrage_idx ON positive_ev(id);
