import nats from "node-nats-streaming";
import { Pool } from "pg";
import { env } from "./env";
import { v4 as uuidV4 } from "uuid";
import { createIdentity } from "./identity";
import { createExpressApp } from "./express";

const stan = nats.connect(env.CLUSTER_ID, `${env.CLIENT_ID}-${uuidV4()}`, {
  url: env.NATS_SERVER,
});

const pool = new Pool({ connectionString: env.DATABASE_URL });
pool.connect();

// ==============================================================
// The pool will emit an error on behalf of any idle clients
// it contains if a backend error or network partition happens
// ==============================================================

pool.on("error", (err) => {
  console.error("Unexpected error on idle client", err);
  process.exit(-1);
});

stan.on("connect", () => {
  console.log("Publisher connected to NATS");

  stan.on("connection_lost", (error) => {
    console.log("disconnected from stan", error);
  });

  stan.on("close", () => {
    console.log("NATS connection closed!");
  });

  const feature = createIdentity(pool, stan);
  const app = createExpressApp(feature, env);

  app.listen(env.PORT, () => {
    console.log(`${env.APP_NAME} started`);
    console.table([
      ["Port", env.PORT],
      ["Environment", env.NODE_ENV],
    ]);
  });
});

process.on("SIGINT", () => {
  console.log("SIGINT");
  stan.close();
});
process.on("SIGTERM", () => {
  console.log("SIGTERM");
  stan.close();
});
