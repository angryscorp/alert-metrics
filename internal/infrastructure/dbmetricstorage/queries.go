package dbmetricstorage

const selectAllMetrics = `
	SELECT * FROM metrics
`

const selectMetric = `
	SELECT value_delta, value_gauge 
	FROM metrics 
	WHERE 
		id = $1 
	  AND 
		type = $2
`

const upsertMetric = `
    INSERT INTO metrics (id, type, value_delta, value_gauge)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (id, type) DO UPDATE SET
		value_delta = CASE 
			WHEN metrics.type = 'counter' 
			THEN metrics.value_delta + EXCLUDED.value_delta
			ELSE EXCLUDED.value_delta
		END,
		value_gauge = EXCLUDED.value_gauge
`

const createTableMetrics = `
	CREATE TABLE IF NOT EXISTS metrics (
		id VARCHAR(255) NOT NULL,
		type VARCHAR(50) NOT NULL,
		value_delta BIGINT,
		value_gauge DOUBLE PRECISION,
		PRIMARY KEY (id, type)
	);
`
