const msg_pb = require("./msg_pb");

export interface Message {
  id: number;
  type: string;
  metadata: {
    traceId: string;
    userId: string;
  };
  seq: number;
  data: any;
  size: number;
  timestamp: number;
  global_position: number;
}

export interface RawMessage {
  id: number;
  seq: string;
  data: Buffer;
  size: number;
  timestamp: string;
  global_position: string;
}

export type ParseMsg = (raw: RawMessage) => Message;

export function deserializeMessage(msg: RawMessage): Message {
  const d = msg_pb.MsgProto.deserializeBinary(msg.data);
  const buff = Buffer.from(d.toObject().data, "base64");
  const text = buff.toString("ascii");
  const { type, metadata, data } = JSON.parse(text);

  console.log("deserializing the data =====================");

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
