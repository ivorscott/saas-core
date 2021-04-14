"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const pg_1 = require("pg");
const env_1 = require("./env");
const msg_1 = require("./msg");
const identity_1 = require("./identity");
const viewdb = new pg_1.Pool({ connectionString: env_1.env.DATABASE_URL }); // view db
const natsdb = new pg_1.Pool({ connectionString: env_1.env.NATS_DB_URL });
viewdb.connect();
viewdb.on("connect", () => console.log("connected to viewdb"));
viewdb.on("error", () => console.log("error occured connecting to viewdb"));
natsdb.connect();
natsdb.on("connect", () => console.log("connected to natsdb"));
viewdb.on("error", () => console.log("error occured connecting to natsdb"));
const messageStore = msg_1.createMessageStore(natsdb);
const identityAggregator = identity_1.createAggregator(viewdb, messageStore);
console.log("starting identity aggregator");
identityAggregator.start();
process.on("SIGINT", () => {
  console.log("SIGINT");
  identityAggregator.stop();
});
process.on("SIGTERM", () => {
  console.log("SIGTERM");
  identityAggregator.stop();
});
