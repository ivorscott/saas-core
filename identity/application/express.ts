import express from "express";
import { Environment } from "./env";
import { Feature } from "./identity";
import { useMiddleware } from "./mid";

interface AppContext {
  feature: Feature;
  env: Environment;
}

function createExpressApp({ feature, env }: AppContext) {
  const app = express();

  useMiddleware(app, env);
  app.use("/v1/users", feature.router);

  return app;
}

export { createExpressApp };
