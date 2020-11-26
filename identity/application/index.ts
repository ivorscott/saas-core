import Knex from "knex";
import { env } from "./env";
import { v4 as uuidV4 } from "uuid";
import { createIdentity } from "./identity";
import { createExpressApp } from "./express";

import nats from "node-nats-streaming";

const client = Promise.resolve(Knex(env.DATABASE_URL));

const stan = nats.connect("devpie-client", uuidV4(), {
  url: "http://nats-svc:4222",
});

stan.on("connect", () => {
  console.log("Publisher connected to NATS");

  stan.on("close", () => {
    console.log("NATS connection closed!");
    process.exit();
  });

  const feature = createIdentity(client, stan);
  const app = createExpressApp(feature, env);

  app.listen(env.PORT, () => {
    console.log(`${env.APP_NAME} started`);
    console.table([
      ["Port", env.PORT],
      ["Environment", env.NODE_ENV],
    ]);
  });
});

process.on("SIGINT", () => stan.close());
process.on("SIGTERM", () => stan.close());
