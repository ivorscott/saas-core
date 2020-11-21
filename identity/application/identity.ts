import express, { Router, Request, Response, NextFunction } from "express";
import Knex from "knex";

interface Queries {
  loadUser: (id: string) => any;
}

interface Handlers {
  findIdentity: (req: Request, res: Response, next: NextFunction) => any;
  saveIdentity: (req: Request, res: Response, next: NextFunction) => void;
}

export interface Feature {
  queries: Queries;
  handlers: Handlers;
  router: Router;
}

function createHandlers({ queries }: { queries: Queries }) {
  function saveIdentity(req: Request, res: Response, next: NextFunction) {
    return res.send({});
  }

  function findIdentity(req: Request, res: Response, next: NextFunction) {
    return res.status(200).send({});
  }

  return {
    findIdentity,
    saveIdentity,
  };
}

function createQueries({ db }: { db: Promise<Knex> }): Queries {
  function loadUser(id: string) {
    return db
      .then((client) => client("users").where({ id }))
      .then((rows) => rows[0]);
  }

  return {
    loadUser,
  };
}

function createIdentity({ db }: { db: Promise<Knex> }): Feature {
  const queries = createQueries({ db });
  const handlers = createHandlers({ queries });

  const router = express.Router();

  router.route("/").post(handlers.saveIdentity);
  router.route("/me").get(handlers.findIdentity);

  return { handlers, queries, router };
}

export { createIdentity };
