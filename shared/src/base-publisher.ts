import { Stan } from "node-nats-streaming";
import { Message } from ".";

export abstract class Publisher<T extends Message> {
  abstract type: T["type"];
  public streamName: string;
  private client: Stan;

  constructor(client: Stan, streamName: string) {
    this.client = client;
    this.streamName = streamName;
  }

  publish(message: T) {
    return new Promise((resolve, reject) => {
      this.client.publish(this.streamName, JSON.stringify(message), (err) => {
        if (err) {
          return reject(err);
        }
        resolve();
        console.log("Message Published", this.type);
      });
    });
  }
}
