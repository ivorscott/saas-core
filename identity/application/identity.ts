import { Stan } from "node-nats-streaming";
import express, { Router, Request, Response } from "express";
import { v4 as uuidV4 } from "uuid";
import { Categories, Commands } from "@devpie/client-events";
import { AddUserPublisher } from "./publish-add-user";
import { Client } from "pg";

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

interface User {
  id: string;
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
  getUser: (id: string) => Promise<User | undefined>;
}

interface Handlers {
  findIdentity: (req: Request, res: Response) => any;
  saveIdentity: (req: Request, res: Response) => any;
}

function createActions(natsClient: Stan, queries: Queries): Actions {
  async function getUser(id: string) {
    return await queries.loadUser(id);
  }

  async function addUser(traceId: string, user: Auth0User) {
    try {
      const userId = uuidV4();
      const streamName = `${Categories.Identity}:command`;
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
      console.log(command);
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
    const auth0User = req.user;
    const user = await actions.getUser(auth0User.auth0Id);
    if (user?.id) {
      return res.status(200).send(user);
    }
    return res.status(404).send({ error: "user not found" });
  }

  async function saveIdentity(req: Request, res: Response) {
    const auth0User = req.body;
    await actions.addUser(req.context.traceId, auth0User);
    return res.status(200).send({});
  }

  return {
    findIdentity,
    saveIdentity,
  };
}

function createQueries(db: Client): Queries {
  function loadUser(auth0Id: string) {
    return db.query(
      "SELECT user_id, auth0_id, email, email_verified, first_name, last_name, picture, created, locale FROM users WHERE auth0_id = $1",
      [auth0Id],
      (err, res) => {
        if (err) throw err;
        db.end();
        return res.rowCount ? res.rows[0] : {};
      },
    );
  }

  return {
    loadUser,
  };
}

export function createIdentity(db: Client, natsClient: Stan): Feature {
  const queries = createQueries(db);
  const actions = createActions(natsClient, queries);
  const handlers = createHandlers(actions);

  const router = express.Router();

  router.route("/").post(handlers.saveIdentity);
  router.route("/me").get(handlers.findIdentity);

  return { handlers, queries, router };
}
