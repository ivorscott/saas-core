import { Stan } from "node-nats-streaming";
import express, { Router, Request, Response } from "express";
import { v4 as uuidV4 } from "uuid";
import { Categories, Commands } from "@devpie/client-events";
import { AddUserPublisher } from "./publish-add-user";
import { Pool } from "pg";
import camelcaseKeys from "camelcase-keys";

export interface Feature {
  actions: Actions;
  queries: Queries;
  handlers: Handlers;
  router: Router;
}

export interface Auth0User {
  auth0Id: string;
  email: string;
  emailVerified: boolean;
  firstName: string;
  lastName: string;
  picture: string;
  locale: string;
}

export interface DBUser {
  user_id: string;
  auth0_id: string;
  email: string;
  email_verified: boolean;
  first_name: string;
  last_name: string;
  picture: string;
  locale: string;
}

export interface User {
  userId: string;
  auth0Id: string;
  email: string;
  emailVerified: boolean;
  firstName: string;
  lastName: string;
  picture: string;
  locale: string;
}

interface Queries {
  loadUser: (auth0Id: string) => Promise<DBUser | undefined>;
}

interface Actions {
  addUser: (traceId: string, user: Auth0User) => void;
  getUser: (id: string) => Promise<User | undefined>;
}

interface Handlers {
  findIdentity: (req: Request, res: Response) => Promise<void>;
  saveIdentity: (req: Request, res: Response) => Promise<void>;
}

export enum ERR {
  UserNotFound,
}

export const errors = [
  { error: "user not found" }, // userNotFound
];

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
  async function getUser(auth0Id: string) {
    return await queries.loadUser(auth0Id).then((record: DBUser) => {
      let user: User;
      if (record) {
        user = (camelcaseKeys(record) as unknown) as User;
      }
      return user;
    });
  }

  // TODO: add Client

  async function addUser(traceId: string, user: Auth0User) {
    try {
      const userId = uuidV4();
      const streamName = `${Categories.Identity}.command`;
      const publisher = new AddUserPublisher(natsClient, streamName);

      const type: Commands.AddUser = Commands.AddUser;

      const command = {
        id: uuidV4(),
        type,
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
    getUser,
  };
}

function createHandlers(actions: Actions) {
  async function findIdentity(req: Request, res: Response) {
    const auth0Id = req.user.sub;
    const user = await actions.getUser(auth0Id);

    if (user) {
      res.status(200).send(user);
      return;
    }
    res.status(404).send(errors[ERR.UserNotFound]);
  }

  async function saveIdentity(req: Request, res: Response) {
    const auth0User = req.body;
    await actions.addUser(req.context.traceId, auth0User);
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
