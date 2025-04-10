DROP TABLE IF EXISTS all_fixtures;

CREATE TABLE IF NOT EXISTS all_fixtures (
    id TEXT,
    start_date TEXT,
    home_team TEXT,
    away_team TEXT,
    sport TEXT,
    league TEXT
) STRICT;

CREATE INDEX IF NOT EXISTS all_fixtures_idx ON all_fixtures(id);
