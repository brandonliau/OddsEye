DROP TABLE IF EXISTS expected_value;

CREATE TABLE IF NOT EXISTS expected_value (
    id TEXT,
    market TEXT,
    grouping_key TEXT,
    selection TEXT,
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

CREATE INDEX IF NOT EXISTS expected_value_idx ON expected_value(id, market, grouping_key);
