import express from "express";
import { Environment } from "./env";
import { Feature } from "./accounting";
import { useMiddleware } from "./mid";

function createExpressApp(feature: Feature, env: Environment) {
  const app = express();

  useMiddleware(app, env);
  app.use("/api/v1/accounting", feature.router);

  return app;
}

export { createExpressApp };
