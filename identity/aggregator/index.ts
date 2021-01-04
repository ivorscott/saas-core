import { Pool } from "pg";
import { env } from "./env";
import { createMessageStore } from "./msg";
import { createAggregator } from "./identity";

const viewdb = new Pool({
  connectionString: env.DATABASE_URL,
  ssl: { rejectUnauthorized: false },
});
const natsdb = new Pool({
  connectionString: env.NATS_DB_URL,
  ssl: { rejectUnauthorized: false },
});

viewdb.connect();
viewdb.on("connect", () => console.log("connected to viewdb"));
viewdb.on("error", () => console.log("error occured connecting to viewdb"));

natsdb.connect();
natsdb.on("connect", () => console.log("connected to natsdb"));
viewdb.on("error", () => console.log("error occured connecting to natsdb"));

const messageStore = createMessageStore(natsdb);
const identityAggregator = createAggregator(viewdb, messageStore);

identityAggregator.start();

process.on("SIGINT", () => {
  console.log("SIGINT");
  identityAggregator.stop();
});
process.on("SIGTERM", () => {
  console.log("SIGTERM");
  identityAggregator.stop();
});
