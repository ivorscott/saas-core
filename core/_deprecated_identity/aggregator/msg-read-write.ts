import { Pool } from "pg";
import { Message, RawMessage } from "./msg";
const message = require("./msg_pb");

const sqlStatements = [
  "SELECT * FROM get_category_stream($1,$2,$3) as (id integer, seq bigint, timestamp bigint, size integer, data bytea, global_position bigint)",
  "SELECT * FROM get_last_read_position($1)",
  "SELECT write_read_position($1,$2)",
];

enum SQL {
  get_category_stream,
  get_last_read_position,
  write_read_position,
}

export function createWrite(db: Pool) {
  return (subscriberStreamName: string, position: number) => {
    if (!position) {
      throw new Error("Message must have a position");
    }
    const values = [subscriberStreamName, position];
    return db.query(sqlStatements[SQL.write_read_position], values);
  };
}

export function createRead(db: Pool) {
  function readLastMessage(subscriberStreamName: string) {
    return db
      .query(sqlStatements[SQL.get_last_read_position], [subscriberStreamName])
      .then((res) => {
        if (!(res.rows[0].length && res.rows[0].seq)) {
          return null;
        }
        return deserializeMessage(res.rows[0]);
      });
  }

  function fetch(streamName: string, projection: any) {
    return readCategory(streamName).then((messages) =>
      project(messages, projection),
    );
  }

  function readCategory(
    categoryStreamName: string,
    fromPosition: number = 0,
    maxMessages: number = 1000,
  ): Promise<Message[]> {
    const values = [categoryStreamName, fromPosition, maxMessages];
    return db
      .query(sqlStatements[SQL.get_category_stream], values)
      .then((res) => {
        return res.rows.map(deserializeMessage);
      });
  }

  return {
    fetch,
    readCategory,
    readLastMessage,
  };
}

// Projects an array of events, running them through a projection
export function project(
  messages: (Message | null)[] | Message[],
  projection: any,
) {
  return;
}

export function deserializeMessage(msg: RawMessage): Message {
  const d = message.MsgProto.deserializeBinary(msg.data);
  const buff = Buffer.from(d.toObject().data, "base64");
  const text = buff.toString("ascii");
  const { type, metadata, data } = JSON.parse(text);

  return {
    id: msg.id,
    seq: parseInt(msg.seq, 10),
    timestamp: parseInt(msg.timestamp, 10),
    size: msg.size,
    global_position: parseInt(msg.global_position, 10),
    type,
    metadata,
    data,
  };
}
