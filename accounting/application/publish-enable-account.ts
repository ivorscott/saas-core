import { Commands, Publisher, EnableAccountingCommand } from "@devpie/client-events";

 export class EnableAccountingPublisher extends Publisher<EnableAccountingCommand> {
  readonly type = Commands.EnableAccounting;
}
