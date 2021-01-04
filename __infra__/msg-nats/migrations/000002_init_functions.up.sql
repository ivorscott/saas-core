--Read all messages by category with global position

CREATE OR REPLACE FUNCTION get_category_stream(category character varying, seq bigint DEFAULT 1, batch_size bigint DEFAULT 1000) RETURNS SETOF record AS $$
WITH category_channels AS (
    SELECT *
    FROM channels
    WHERE name LIKE '%' || category || '.________-____-____-____-____________' || '%'
), all_messages AS (
    SELECT *, ROW_NUMBER() OVER (ORDER BY (timestamp)) AS global_position
    FROM messages
)

SELECT *
FROM all_messages
WHERE id IN (SELECT id FROM category_channels) AND global_position >= seq
ORDER BY timestamp
LIMIT batch_size
$$ LANGUAGE SQL;

--Get last read position message

CREATE OR REPLACE FUNCTION get_last_read_position(sub_id varchar (1024)) RETURNS sub_sequence AS $$
SELECT *
FROM sub_sequence
WHERE subid=$1
ORDER BY timestamp DESC
LIMIT 1
$$ LANGUAGE SQL;

--Write last position message

CREATE OR REPLACE FUNCTION write_read_position(streamName varchar(1024), seq integer) RETURNS void AS $$
INSERT INTO sub_sequence (
    subid,
    seq
) VALUES (
    streamName,
    seq
)
ON CONFLICT DO NOTHING
$$ LANGUAGE SQL;