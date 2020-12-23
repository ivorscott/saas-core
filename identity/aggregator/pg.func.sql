CREATE OR REPLACE FUNCTION readCategoryStream(category character varying) RETURNS SETOF channels AS $$
    SELECT *
    FROM channels
    WHERE name LIKE '%' || category || '.________-____-____-____-____________' || '%'
$$ LANGUAGE SQL;

-- CREATE OR REPLACE FUNCTION recoverChannelMsgs(channelID integer) RETURNS SETOF messages AS $$
--     SELECT COUNT(seq),
--     COALESCE(MIN(seq), 0),
--     COALESCE(MAX(seq), 0),
--     COALESCE(SUM(size), 0),
--     COALESCE(MAX(timestamp), 0)
--     FROM Messages WHERE id=channelID
-- $$ LANGUAGE SQL; (sqlRecoverChannelMsgs)



-- "SELECT timestamp, data FROM Messages WHERE id=? AND seq=?", (sqlLookupMsg)