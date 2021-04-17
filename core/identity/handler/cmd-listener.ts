import { Message as MessageUtils } from "node-nats-streaming";
import {
  Message,
  Events,
  Commands,
  Categories,
  UserAddedEvent,
  CommandListener,
} from "@devpie/client-events";
import { createMessageStore } from "./msg-store";
import { v4 as uuidV4 } from "uuid";
import { UserAddedPublisher } from "./publish-user-added";

export class IdentityCommandListener extends CommandListener {
  readonly category = Categories.Identity;
  queueGroupName = "mh-identity";

  async onMessage(command: Message, utils: MessageUtils) {
    switch (command.type) {
      case Commands.AddUser:
        const { fetch } = createMessageStore(this.db);
        const entityStream = `${this.category}.${command.metadata.userId}`;
        const events = await fetch(entityStream);

        if (events.length === 0) {
          const publisher = new UserAddedPublisher(this.client, entityStream);
          const event: UserAddedEvent = {
            id: uuidV4(),
            type: Events.UserAdded,
            metadata: command.metadata,
            data: command.data,
          };
          publisher.publish(event);
        }
        utils.ack();
    }
  }
}
