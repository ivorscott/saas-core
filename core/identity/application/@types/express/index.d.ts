import * as express from "express";

declare global {
  namespace Express {
    interface Request {
      user: any;
      context: {
        userId: string;
        traceId: string;
      };
    }
  }
}
