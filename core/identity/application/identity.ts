import { Pool } from "pg";
import { Stan } from "node-nats-streaming";
import { DBUser, ReqContext, Auth0User, User } from "types";
import { Categories, Commands } from "@devpie/client-events";
import { AddUserPublisher } from "./publish-add-user";
import { v4 as uuidV4 } from "uuid";
import express, { Router, Request, Response } from "express";
import camelcaseKeys from "camelcase-keys";

export interface Feature {
  actions: Actions;
  queries: Queries;
  handlers: Handlers;
  router: Router;
}

interface Queries {
  loadUser: (auth0Id: string) => Promise<DBUser | undefined>;
}

interface Actions {
  addUser: (ctx: ReqContext, user: Auth0User) => void;
  getUser: (ctx: ReqContext, id: string) => Promise<User | undefined>;
}

interface Handlers {
  findIdentity: (req: Request, res: Response) => Promise<void>;
  saveIdentity: (req: Request, res: Response) => Promise<void>;
}

export enum ERR {
  UserNotFound,
}

export const errors = [{ error: "user not found" }];

export enum SQL {
  GetUser,
  AddUser,
}

export const sqlStatements = [
  "SELECT * FROM users WHERE auth0_id = $1 LIMIT 1", // getUser
];

export function createQueries(db: Pool): Queries {
  function loadUser(auth0Id: string): Promise<DBUser | undefined> {
    return db
      .query(sqlStatements[SQL.GetUser], [auth0Id])
      .then((res) => res.rows[0]);
  }

  return {
    loadUser,
  };
}

export function createActions(natsClient: Stan, queries: Queries): Actions {
  async function getUser(ctx: ReqContext, auth0Id: string) {
    return await queries.loadUser(auth0Id).then((record: DBUser) => {
      let user: User;
      if (record) {
        user = (camelcaseKeys(record) as unknown) as User;
      }
      return user;
    });
  }

  async function addUser(ctx: ReqContext, user: Auth0User) {
    try {
      const cmdStream = `${Categories.Identity}.command`;
      const publisher = new AddUserPublisher(natsClient, cmdStream);

      const type: Commands.AddUser = Commands.AddUser;

      const command = {
        id: uuidV4(),
        type,
        metadata: ctx,
        data: { id: uuidV4(), ...user },
      };

      console.log("AddUser: ", command);
      await publisher.publish(command);
    } catch (error) {
      console.log("error:", error);
    }
    return;
  }

  return {
    addUser,
    getUser,
  };
}

function createHandlers(actions: Actions) {
  async function findIdentity(req: Request, res: Response) {
    const auth0Id = req.user.sub;
    const user = await actions.getUser(req.context, auth0Id);

    if (user) {
      res.status(200).send(user);
      return;
    }
    res.status(404).send(errors[ERR.UserNotFound]);
  }

  async function saveIdentity(req: Request, res: Response) {
    const auth0User = req.body as Auth0User;
    await actions.addUser(req.context, auth0User);
    res.status(200).end();
  }

  return {
    findIdentity,
    saveIdentity,
  };
}

export function createIdentity(
  db: Pool,
  natsClient: Stan,
): {
  router: Router;
  handlers: {
    findIdentity: (req: Request, res: Response) => Promise<void>;
    saveIdentity: (req: Request, res: Response) => Promise<void>;
  };
  actions: Actions;
  queries: Queries;
} {
  const queries = createQueries(db);
  const actions = createActions(natsClient, queries);
  const handlers = createHandlers(actions);

  const router = express.Router();

  router.route("/").post(handlers.saveIdentity);
  router.route("/me").get(handlers.findIdentity);

  return { handlers, actions, queries, router };
}
