import { deserializeMessage, ParseMsg } from "./msg-deserialize";
import { Pool } from "pg";

enum SQL {
  GetEntityStream,
}

const sqlStatements = [
  `SELECT * FROM get_stream($1,$2,$3) as (id integer, seq bigint, timestamp bigint, size integer, data bytea, global_position bigint)`,
];

function createRead(db: Pool, parser: ParseMsg) {
  async function fetch(
    streamName: string,
    fromPosition = 0,
    maxMessages = 1000,
  ) {
    const values = [streamName, fromPosition, maxMessages];
    return await db
      .query(sqlStatements[SQL.GetEntityStream], values)
      .then((res) => {
        return res.rows.map(parser);
      });
  }

  return { fetch };
}

export function createMessageStore(db: Pool) {
  db.connect();

  const { fetch } = createRead(db, deserializeMessage);

  return {
    fetch,
  };
}
