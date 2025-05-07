CREATE TABLE IF NOT EXISTS metrics (
    id VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    value_delta BIGINT,
    value_gauge DOUBLE PRECISION,
    PRIMARY KEY (id, type)
);
