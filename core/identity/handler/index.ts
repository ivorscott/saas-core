import { Categories } from "@devpie/client-events";
import { reconnectHandler, natsInstance as natsClient } from "./nats-streaming";
import { AddUserListener } from "./listen-add-user";

reconnectHandler();

const streamName = `${Categories.Identity}.command`;
new AddUserListener(natsClient, streamName).listen();
