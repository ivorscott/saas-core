import { Commands, Publisher, AddUserCommand } from "@devpie/client-events";

export class AddUserPublisher extends Publisher<AddUserCommand> {
  readonly subject = Commands.AddUser;
}
