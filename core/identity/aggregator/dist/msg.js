"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createMessageStore = void 0;
const msg_subscribe_1 = require("./msg-subscribe");
const msg_read_write_1 = require("./msg-read-write");
function createMessageStore(db) {
  const write = msg_read_write_1.createWrite(db);
  const read = msg_read_write_1.createRead(db);
  const createSubscription = msg_subscribe_1.configureCreateSubscription(
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
exports.createMessageStore = createMessageStore;
