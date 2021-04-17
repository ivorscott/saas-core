import { v4 } from "uuid";
import { Request, Response, NextFunction } from "express";

function requestContext(req: Request, _res: Response, next: NextFunction) {
  req.context = {
    traceId: v4(),
    userId: req.user.sub,
  };
  next();
}

export { requestContext };
