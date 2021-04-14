import { Pool } from "pg";
import { env } from "./env";
import { createMessageStore } from "./msg";
import { createAggregator } from "./identity";

let viewClient: Pool
let natsClient: Pool
let identityAggregator: { start: () => void; stop: () => void; }

const viewdb = new Pool({
  connectionString: env.DATABASE_URL,
  ssl: { rejectUnauthorized: false },
});

const natsdb = new Pool({
  connectionString: env.NATS_DB_URL,
  ssl: { rejectUnauthorized: false },
});

function main() {
  Promise.all([viewdb.connect,natsdb.connect]).then(() => {
    const messageStore = createMessageStore(natsdb);
    identityAggregator = createAggregator(viewdb, messageStore);
    identityAggregator.start();
  })
}

const reconnectHandler = ()=> {
  console.log("[RECONNECTING]")
  identityAggregator.stop()
  Promise.all([viewClient.end,natsClient.end]).then(()=> {
    main()
  })
}

viewdb.on("connect", () => console.log("[CONNECTED] to viewdb"));
viewdb.on("error", () => console.log("[ERROR] connecting to viewdb"));
natsdb.on("connect", () => console.log("[CONNECTED] to natsdb"));
natsdb.on("error", () => console.log("[ERROR] connecting to natsdb"));

process.on("SIGINT", () => {
  console.log("[SIGINT]");
   reconnectHandler()
});

process.on("SIGTERM",  () => {
  console.log("[SIGTERM]");
   reconnectHandler()
});

main()
