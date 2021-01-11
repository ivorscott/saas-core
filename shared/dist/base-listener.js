"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Listener = void 0;
class Listener {
    constructor(client, streamName) {
        this.ackWait = 5 * 1000; // 5000 milliseconds
        this.client = client;
        this.streamName = streamName;
    }
    subscriptionOptions() {
        return this.client
            .subscriptionOptions()
            .setDeliverAllAvailable()
            .setManualAckMode(true)
            .setAckWait(this.ackWait)
            .setDurableName(this.queueGroupName);
    }
    listen() {
        const subscription = this.client.subscribe(this.type, this.queueGroupName, this.subscriptionOptions());
        subscription.on("message", (msg) => {
            console.log(`Message received: ${this.streamName} / ${this.queueGroupName}`);
            const parsedData = this.parseMessage(msg);
            this.onMessage(parsedData, msg);
        });
    }
    parseMessage(msg) {
        const data = msg.getData();
        return typeof data === "string"
            ? JSON.parse(data)
            : JSON.parse(data.toString("utf-8")); // parse buffer
    }
}
exports.Listener = Listener;
