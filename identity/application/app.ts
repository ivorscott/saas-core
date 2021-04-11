import express from "express";
import { Environment } from "./env";
import { Feature } from "./identity";
import { useMiddleware } from "./mid";

function createExpressApp(feature: Feature, env: Environment) {
  const app = express();

  useMiddleware(app, env);
  app.use("/api/v1/users", feature.router);

  return app;
}

export { createExpressApp };
