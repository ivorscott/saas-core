"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Publisher = void 0;
class Publisher {
    constructor(client, streamName) {
        this.client = client;
        this.streamName = streamName;
    }
    publish(message) {
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
exports.Publisher = Publisher;
