import { Events, Publisher, UserAddedEvent } from "@devpie/client-events";

export class UserAddedPublisher extends Publisher<UserAddedEvent> {
  readonly type = Events.UserAdded;
}
