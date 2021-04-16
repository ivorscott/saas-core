import { Pool } from "pg";
import { env } from "./env";
import { v4 as uuidV4 } from "uuid";
import { Message } from "node-nats-streaming";
import { UserAddedPublisher } from "./publish-user-added";
import { deserializeMessage } from "./msg-deserialize";
import { createMessageStore } from "./msg";

import {
  Events,
  Commands,
  Listener,
  AddUserCommand,
  UserAddedEvent,
} from "@devpie/client-events";

export class AddUserListener extends Listener<AddUserCommand> {
  type: Commands.AddUser;
  queueGroupName = "mh-identity";

  async onMessage(command: AddUserCommand, msg: Message) {
    const db = new Pool({
      connectionString: env.NATS_DB_URL,
      ssl: { rejectUnauthorized: false },
    });

    console.log("Test");

    const { replayEvents } = createMessageStore(db, deserializeMessage);

    // check events
    const events = await replayEvents(this.streamName);

    console.log("Historical events: ", events);

    const result = events.map(
      (e: any) => e.data.auth0Id === command.data.auth0Id,
    );

    // Never seen it
    if (result.length === 0) {
      const publisher = new UserAddedPublisher(this.client, this.streamName);

      const event: UserAddedEvent = {
        id: uuidV4(),
        type: Events.UserAdded,
        metadata: {
          traceId: uuidV4(),
          userId: command.data.id,
        },
        data: command.data,
      };

      console.log("New event: ", event);

      publisher.publish(event);
    }
    msg.ack();
  }
}
