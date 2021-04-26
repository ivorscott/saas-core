"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createMessageStore = void 0;
const subscribe_1 = require("./subscribe");
const read_write_1 = require("./read-write");
function createMessageStore(db) {
  console.log("creating message store");
  const write = read_write_1.createWrite(db);
  const read = read_write_1.createRead(db);
  const createSubscription = subscribe_1.configureCreateSubscription(
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
