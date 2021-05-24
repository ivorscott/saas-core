CREATE OR REPLACE FUNCTION get_stream(stream_name character varying, seq bigint DEFAULT 1, batch_size bigint DEFAULT 1000) RETURNS SETOF record AS $$
WITH channel AS (
    SELECT *
    FROM channels
    WHERE name = stream_name
), all_messages AS (
    SELECT *, ROW_NUMBER() OVER (ORDER BY (timestamp)) AS global_position
    FROM messages
)

SELECT *
FROM all_messages
WHERE id IN (SELECT id FROM channel) AND global_position >= seq
ORDER BY timestamp
LIMIT batch_size
$$ LANGUAGE SQL;
