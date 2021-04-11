import { env } from "./env";
import { reconnectHandler, natsInstance } from "./nats-streaming";
import { createAccounting } from "./accounting";
import { createExpressApp } from "./app";

reconnectHandler();

const feature = createAccounting(natsInstance);
const app = createExpressApp(feature, env);

app.listen(env.PORT, () => {
  console.log(`${env.APP_NAME} started.`);
  console.table([
    ["Port", env.PORT],
    ["Environment", env.NODE_ENV],
  ]);
});
