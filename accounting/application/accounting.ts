import { Stan } from "node-nats-streaming";
import express, { Router, Request, Response } from "express";
import { v4 as uuidV4 } from "uuid";
import { Categories, Commands } from "@devpie/client-events";
import { EnableAccountingPublisher } from "./publish-enable-account";

export interface Feature {
  actions: Actions;
  handlers: Handlers;
  router: Router;
}

export interface NewAccount {
  auth0Id: string;
  token: string;
}

interface Actions {
  addAccount: (traceId: string, user: NewAccount) => void;
}

interface Handlers {
  saveAccount: (req: Request, res: Response) => Promise<void>;
}

export function createActions(natsClient: Stan): Actions {

  async function addAccount(traceId: string, account: NewAccount) {
    try {
      const userId = uuidV4();
      const streamName = `${Categories.Accounting}.command`;
      const publisher = new EnableAccountingPublisher(natsClient, streamName);

      const type: Commands.EnableAccounting = Commands.EnableAccounting;

      const command = {
        id: uuidV4(),
        type,
        metadata: {
          traceId,
          userId,
        },
        data: { id: userId, ...account },
      };
      await publisher.publish(command);
    } catch (error) {
      console.log("error:", error);
    }
    return;
  }

  return {
    addAccount,
  };
}

function createHandlers(actions: Actions) {
  async function saveAccount(req: Request, res: Response) {
    const account = req.body;
    await actions.addAccount(req.context.traceId, account);
    res.status(200).end();
  }

  return {
    saveAccount,
  };
}

export function createAccounting(
  natsClient: Stan,
): {
  router: Router;
  handlers: {
    saveAccount: (req: Request, res: Response) => Promise<void>;
  };
  actions: Actions;
} {
  const actions = createActions(natsClient);
  const handlers = createHandlers(actions);

  const router = express.Router();

  router.route("/").post(handlers.saveAccount);

  return { handlers, actions, router };
}
