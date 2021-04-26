import Bluebird from "bluebird";
import { Message } from "./msg";
import { QueryResult } from "pg";

export function configureCreateSubscription(
  readCategory: (
    categoryStreamName: string,
    fromPosition?: number,
    maxMessages?: number,
  ) => Promise<Message[]>,
  readLastMessage: (subscriberStreamName: string) => Promise<Message | null>,
  write: (
    subscriberStreamName: string,
    position: number,
  ) => Promise<QueryResult>,
) {
  return (
    streamName: string,
    handlers: any,
    subscriberId: string,
    messagesPerTick = 100,
    positionUpdateInterval: number = 100,
    tickIntervalsMs: number = 100,
  ) => {
    let currentPosition = 0;
    let messagesSinceLastPositionWrite = 0;
    let continuePolling = true;

    function writePosition(position: number) {
      return write(subscriberId, position);
    }

    function loadPosition() {
      readLastMessage(subscriberId).then((message: Message | null) => {
        currentPosition = message ? message.seq : 0;
      });
    }

    function updateReadPosition(position: number) {
      currentPosition = position;
      messagesSinceLastPositionWrite += 1;

      if (messagesSinceLastPositionWrite === positionUpdateInterval) {
        messagesSinceLastPositionWrite = 0;
        return writePosition(position);
      }

      return Promise.resolve(true);
    }

    function getNextBatchOfMessages() {
      return readCategory(streamName, currentPosition + 1, messagesPerTick);
    }

    function handleMessage(message: Message) {
      const handler = handlers[message.type];
      return handler ? handler(message) : Promise.resolve(true);
    }

    function processBatch(messages: Message[]) {
      return Bluebird.each(messages, (message) =>
        handleMessage(message)
          .then(() => updateReadPosition(message.global_position))
          .catch((err: Error) => {
            logError(message, err);

            throw err;
          }),
      ).then(() => messages.length);
    }

    function logError(lastMessage: Message, error: Error) {
      console.error(`error processing: ${streamName}`, lastMessage, error);
    }

    function tick() {
      return getNextBatchOfMessages()
        .then(processBatch)
        .catch((err: Error) => {
          console.error("Error processing batch", err);
          stop();
        });
    }

    async function poll() {
      await loadPosition();

      while (continuePolling) {
        const messagesProcessed = await tick();
        if (messagesProcessed === 0) {
          await Bluebird.delay(tickIntervalsMs);
        }
      }
    }

    function start() {
      console.log(`Starting ${subscriberId}`);
      return poll();
    }

    function stop() {
      console.log(`Stopped ${subscriberId}`);
      continuePolling = false;
    }

    // return functions for testing purposes
    return {
      loadPosition,
      start,
      stop,
      writePosition,
    };
  };
}
