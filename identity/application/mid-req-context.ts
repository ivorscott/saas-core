import { v4 } from "uuid";
import { Request, Response, NextFunction } from "express";

function requestContext(req: Request, _res: Response, next: NextFunction) {
  console.log(req.user);
  req.context = {
    traceId: v4(),
    userId: "", // use cookie-session
  };
  next();
}

export { requestContext };
