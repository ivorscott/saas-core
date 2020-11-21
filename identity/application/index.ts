import Knex from "knex";
import { env } from "./env";
import { createIdentity } from "./identity";
import { createExpressApp } from "./express";

const client = Promise.resolve(Knex(env.DATABASE_URL));

// integrate NATS Streaming

const feature = createIdentity({ db: client });

const app = createExpressApp({ feature, env });

app.listen(env.PORT, () => {
  console.log(`${env.APP_NAME} started`);
  console.table([
    ["Port", env.PORT],
    ["Environment", env.NODE_ENV],
  ]);
});
