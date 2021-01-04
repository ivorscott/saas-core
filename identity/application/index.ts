import nats, {Stan} from "node-nats-streaming";
import { Pool } from "pg";
import { env } from "./env";
import { v4 as uuidV4 } from "uuid";
import { createIdentity } from "./identity";
import { createExpressApp } from "./app";

let stan: Stan = natsConnect()

const viewDB = new Pool({ connectionString: env.DATABASE_URL });
viewDB.connect();

viewDB.on("error", (err) => {
  console.error("Unexpected error on idle client", err);
  process.exit(1);
});

stan.on("connection_lost", (error) => {
  console.log("Disconnected from stan.", error);
});

stan.on("close", () => {
  console.log("NATS Streaming connection closed! Reconnecting...");
  stan = natsConnect();
});

stan.on("connect", () => {
  console.log("Publisher connected to NATS");

  const feature = createIdentity(viewDB, stan);
  const app = createExpressApp(feature, env);

  app.listen(env.PORT, () => {
    console.log(`${env.APP_NAME} started.`);
    console.table([
      ["Port", env.PORT],
      ["Environment", env.NODE_ENV],
    ]);
  });
});

process.on("SIGINT", () => {
  console.log("SIGINT detected.");
  stan.close();
});

process.on("SIGTERM", () => {
  console.log("SIGTERM detected.");
  stan.close();
});

function natsConnect(): Stan {
  return nats.connect(env.CLUSTER_ID, `${env.CLIENT_ID}-${uuidV4()}`, {
    url: env.NATS_SERVER,
  });
}