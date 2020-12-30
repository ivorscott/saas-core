"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deserializeMessage = exports.project = exports.createRead = exports.createWrite = void 0;
const message = require("./msg_pb");
const sqlStatements = [
  "SELECT * FROM get_category_stream($1,$2,$3) as (id integer, seq bigint, timestamp bigint, size integer, data bytea, global_position bigint)",
  "SELECT * FROM get_last_read_position($1)",
  "SELECT write_read_position($1,$2)",
];
var SQL;
(function (SQL) {
  SQL[(SQL["get_category_stream"] = 0)] = "get_category_stream";
  SQL[(SQL["get_last_read_position"] = 1)] = "get_last_read_position";
  SQL[(SQL["write_read_position"] = 2)] = "write_read_position";
})(SQL || (SQL = {}));
function createWrite(db) {
  return (subscriberStreamName, position) => {
    if (!position) {
      throw new Error("Message must have a position");
    }
    const values = [subscriberStreamName, position];
    return db.query(sqlStatements[SQL.write_read_position], values);
  };
}
exports.createWrite = createWrite;
function createRead(db) {
  function readLastMessage(subscriberStreamName) {
    return db
      .query(sqlStatements[SQL.get_last_read_position], [subscriberStreamName])
      .then((res) => {
        if (!(res.rows[0].length && res.rows[0].seq)) {
          return null;
        }
        return deserializeMessage(res.rows[0]);
      });
  }
  function fetch(streamName, projection) {
    return readCategory(streamName).then((messages) =>
      project(messages, projection),
    );
  }
  function readCategory(
    categoryStreamName,
    fromPosition = 0,
    maxMessages = 1000,
  ) {
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
exports.createRead = createRead;
// Projects an array of events, running them through a projection
function project(messages, projection) {
  return;
}
exports.project = project;
function deserializeMessage(msg) {
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
exports.deserializeMessage = deserializeMessage;
