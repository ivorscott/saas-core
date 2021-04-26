import bodyParser from "body-parser";
import cors from "cors";
import { checkJwt } from "./mid-auth";
import { requestContext } from "./mid-req-context";
import { Application } from "express";
import { serverError } from "./mid-error-handler";
import { Environment } from "./env";

function useMiddleware(app: Application, env: Environment) {
  app.use(bodyParser.json());
  app.use(cors());
  app.use(checkJwt);
  app.use(requestContext);
  app.use(serverError);
}

export { useMiddleware };
