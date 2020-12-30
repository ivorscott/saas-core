"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.project = exports.createRead = exports.createWrite = exports.sqlStatements = exports.SQL = void 0;
const message_1 = require("./message");
var SQL;
(function (SQL) {
  SQL[(SQL["get_category_stream"] = 0)] = "get_category_stream";
  SQL[(SQL["get_last_read_position"] = 1)] = "get_last_read_position";
  SQL[(SQL["write_read_position"] = 2)] = "write_read_position";
})((SQL = exports.SQL || (exports.SQL = {})));
exports.sqlStatements = [
  "SELECT * FROM get_category_stream($1,$2,$3) as (id integer, seq bigint, timestamp bigint, size integer, data bytea, global_position bigint)",
  "SELECT * FROM get_last_read_position($1)",
  "SELECT write_read_position($1,$2,$3)",
];
function createWrite(db) {
  return (subscriberStreamName, position) => {
    if (!position) {
      throw new Error("Message must have a position");
    }
    const values = [subscriberStreamName, position];
    console.log("writing position");
    return db.query(exports.sqlStatements[SQL.write_read_position], values);
  };
}
exports.createWrite = createWrite;
function createRead(db) {
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
    const query = exports.sqlStatements[SQL.get_category_stream];
    console.log("printing query.....", query);
    const values = [categoryStreamName, fromPosition, maxMessages];
    console.log("printing query values.....", values);
    return db.query(query, values).then((res) => {
      console.log("testing..... read category", res);
      console.log(
        "testing query results...",
        res.rows.length,
        JSON.stringify(res.rows),
      );
      return res.rows.map(message_1.deserializeMessage);
    });
  }
  function readLastMessage(subscriberStreamName) {
    return db
      .query(exports.sqlStatements[SQL.get_last_read_position], [
        subscriberStreamName,
      ])
      .then((res) => message_1.deserializeMessage(res.rows[0]));
  }
  return {
    fetch,
    readCategory,
    readLastMessage,
  };
}
exports.createRead = createRead;
// Projects an array of events, running them through a projection
function project(events, projection) {
  return;
}
exports.project = project;
