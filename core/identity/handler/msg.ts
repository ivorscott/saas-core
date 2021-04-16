import { Message } from "@devpie/client-events";
import { ParseMsg } from "./msg-deserialize";
import { Pool } from "pg";

const sqlStatements = [
  "SELECT * FROM get_category_stream($1,$2,$3) as (id integer, seq bigint, timestamp bigint, size integer, data bytea, global_position bigint)",
];

enum SQL {
  get_category_stream,
}

function createReplayEvents(db: Pool, parser: ParseMsg) {
  async function replayEvents(
    streamName: string,
    fromPosition = 0,
    maxMessages = 1000,
  ) {
    const values = [streamName, fromPosition, maxMessages];
    return await db
      .query(sqlStatements[SQL.get_category_stream], values)
      .then((res) => {
        return res.rows.map(parser);
      });
  }

  return replayEvents;
}

export function createMessageStore(db: Pool, parser: ParseMsg) {
  db.connect();

  const replayEvents = createReplayEvents(db, parser);

  return {
    replayEvents,
  };
}
