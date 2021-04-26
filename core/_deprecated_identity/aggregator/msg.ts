import { configureCreateSubscription } from "./msg-subscribe";
import { createRead, createWrite } from "./msg-read-write";
import { Pool } from "pg";

export interface Message {
  id: number;
  type: string;
  metadata: {
    traceId: string;
    userId: string;
  };
  seq: number;
  data: MessageData;
  size: number;
  timestamp: number;
  global_position: number;
}

export interface MessageData {
  auth0Id: string;
  email: string;
  emailVerified: boolean;
  firstName: string;
  id: string;
  lastName: string;
  locale: string;
  picture: string;
}

export interface RawMessage {
  id: number;
  seq: string;
  data: Buffer;
  size: number;
  timestamp: string;
  global_position: string;
}

export function createMessageStore(db: Pool) {
  const write = createWrite(db);
  const read = createRead(db);
  const createSubscription = configureCreateSubscription(
    read.readCategory,
    read.readLastMessage,
    write,
  );

  return {
    write,
    createSubscription,
    readCategory: read.readCategory,
    readLastMessage: read.readLastMessage,
    fetch: read.fetch,
    stop: db.end,
  };
}
