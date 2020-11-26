import Knex from "knex";
import { Stan } from "node-nats-streaming";
import express, { Router, Request, Response, NextFunction } from "express";
import { v4 as uuidV4 } from "uuid";
import { Commands } from "@devpie/client-events";
import { AddUserPublisher } from "./publish-add-user";

export interface Feature {
  queries: Queries;
  handlers: Handlers;
  router: Router;
}

interface Auth0User {
  auth0Id: string;
  email: string;
  emailVerified: boolean;
  firstName: string;
  lastName: string;
  picture: string;
  locale: string;
}

interface Queries {
  loadUser: (id: string) => any;
}

interface Actions {
  addUser: (traceId: string, user: Auth0User) => void;
}

interface Handlers {
  findIdentityIfExists: (
    req: Request,
    res: Response,
    next: NextFunction,
  ) => any;
}

function createActions(natsClient: Stan, queries: Queries): Actions {
  async function addUser(traceId: string, user: Auth0User) {
    try {
      const publisher = new AddUserPublisher(natsClient);
      const userId = uuidV4();

      const subject: Commands.AddUser = Commands.AddUser;

      const command = {
        id: uuidV4(),
        subject,
        metadata: {
          traceId,
          userId,
        },
        data: { id: userId, ...user },
      };

      await publisher.publish(command);
    } catch (error) {
      console.log("error:", error);
    }
    return;
  }

  return {
    addUser,
  };
}

function createHandlers(actions: Actions) {
  function findIdentityIfExists(req: Request, res: Response) {
    const auth0User = req.user;
    actions.addUser(req.context.traceId, auth0User);
    return res.status(200).send({});
    // .send({ error: "user not found" });
  }

  return {
    findIdentityIfExists,
  };
}

function createQueries(db: Promise<Knex>): Queries {
  function loadUser(id: string) {
    return db
      .then((client) => client("users").where({ id }))
      .then((rows) => rows[0]);
  }

  return {
    loadUser,
  };
}

function createIdentity(db: Promise<Knex>, natsClient: Stan): Feature {
  const queries = createQueries(db);
  const actions = createActions(natsClient, queries);
  const handlers = createHandlers(actions);

  const router = express.Router();

  router.route("/me").get(handlers.findIdentityIfExists);

  return { handlers, queries, router };
}

export { createIdentity };
