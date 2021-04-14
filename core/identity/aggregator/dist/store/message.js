"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.deserializeMessage = void 0;
const message = require("./message_pb");
function deserializeMessage(msg) {
  console.log(msg);
  if (!msg || !msg.data) {
    return null;
  }
  const d = message.MsgProto.deserializeBinary(msg.data);
  const buff = new Buffer(d.toObject().data, "base64");
  const text = buff.toString("ascii");
  const { type, metadata, data } = JSON.parse(text);
  return {
    id: msg.id,
    seq: parseInt(msg.seq, 10),
    timestamp: msg.timestamp,
    size: msg.size,
    global_position: parseInt(msg.global_position, 10),
    type,
    metadata,
    data,
  };
}
exports.deserializeMessage = deserializeMessage;
