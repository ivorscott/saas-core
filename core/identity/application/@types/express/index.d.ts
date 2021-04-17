import * as express from "express";
import { Auth0Claims, ReqContext } from "types";

declare global {
  namespace Express {
    interface Request {
      user: Auth0Claims;
      context: ReqContext;
    }
  }
}
