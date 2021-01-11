import { Stan } from "node-nats-streaming";
import { Message } from ".";
export declare abstract class Publisher<T extends Message> {
    abstract type: T["type"];
    streamName: string;
    private client;
    constructor(client: Stan, streamName: string);
    publish(message: T): Promise<unknown>;
}
