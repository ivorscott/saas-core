import { Pool } from "pg";
import { env } from "./env";
import { reconnectHandler, natsInstance } from "./nats-streaming";
import { createIdentity } from "./identity";
import { createExpressApp } from "./app";

const viewDB = new Pool({ connectionString: env.DATABASE_URL });
viewDB.connect();

viewDB.on("error", (err) => {
  console.error("Unexpected error on idle client", err);
  process.exit(1);
});

reconnectHandler();

const feature = createIdentity(viewDB, natsInstance);
const app = createExpressApp(feature, env);

app.listen(env.PORT, () => {
  console.log(`${env.APP_NAME} started.`);
  console.table([
    ["Port", env.PORT],
    ["Environment", env.NODE_ENV],
  ]);
});
