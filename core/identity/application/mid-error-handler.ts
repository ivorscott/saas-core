import { NextFunction, Request, Response } from "express";

export class HttpException extends Error {
  public status: number;
  public message: string;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.message = message;
  }
}

function serverError(
  err: HttpException,
  req: Request,
  res: Response,
  _next: NextFunction,
) {
  const traceId = req.context ? req.context.traceId : "";

  console.log(traceId, err);
  res.status(500).send("error");
}
export { serverError };
