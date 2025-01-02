DROP TABLE IF EXISTS all_fixtures;
DROP TABLE IF EXISTS all_odds;
DROP TABLE IF EXISTS fair_odds;
DROP TABLE IF EXISTS arbitrage;
DROP TABLE IF EXISTS positive_ev;

DROP TABLE IF EXISTS all_odds_arbitrage;
DROP TABLE IF EXISTS all_odds_positive_ev;

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
    norm_selection TEXT,
    selection_line TEXT,
    points REAL,
    sportsbook TEXT,
    price REAL,
    url TEXT,
    grouping_key TEXT
) STRICT;

CREATE TABLE IF NOT EXISTS fair_odds (
    id TEXT,
    market TEXT,
    selection TEXT,
    norm_selection TEXT,
    selection_line TEXT,
    points REAL,
    sportsbook TEXT,
    price REAL,
    url TEXT,
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
    sportsbook_α TEXT, 
    sportsbook_β TEXT, 
    sportsbook_γ TEXT, 
    total_implied_probability REAL,
    vig REAL
) STRICT;

CREATE TABLE IF NOT EXISTS positive_ev (
    id TEXT,
    market TEXT,
    sportsbook TEXT,
    selection TEXT,
    price REAL,
    fair_odds REAL,
    expected_value REAL
) STRICT;

CREATE INDEX IF NOT EXISTS id_all_fixtures_idx ON all_fixtures(id);
CREATE INDEX IF NOT EXISTS id_market_all_odds_idx ON all_odds(id, market);
CREATE INDEX IF NOT EXISTS id_fair_odds_idx ON fair_odds(id, market);
CREATE INDEX IF NOT EXISTS id_arbitrage_idx ON arbitrage(id);
CREATE INDEX IF NOT EXISTS id_arbitrage_idx ON positive_ev(id);
