"use strict";
var __awaiter =
  (this && this.__awaiter) ||
  function (thisArg, _arguments, P, generator) {
    function adopt(value) {
      return value instanceof P
        ? value
        : new P(function (resolve) {
            resolve(value);
          });
    }
    return new (P || (P = Promise))(function (resolve, reject) {
      function fulfilled(value) {
        try {
          step(generator.next(value));
        } catch (e) {
          reject(e);
        }
      }
      function rejected(value) {
        try {
          step(generator["throw"](value));
        } catch (e) {
          reject(e);
        }
      }
      function step(result) {
        result.done
          ? resolve(result.value)
          : adopt(result.value).then(fulfilled, rejected);
      }
      step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
  };
var __importDefault =
  (this && this.__importDefault) ||
  function (mod) {
    return mod && mod.__esModule ? mod : { default: mod };
  };
Object.defineProperty(exports, "__esModule", { value: true });
exports.configureCreateSubscription = void 0;
const bluebird_1 = __importDefault(require("bluebird"));
function configureCreateSubscription(readCategory, readLastMessage, write) {
  return (
    streamName,
    handlers,
    subscriberId,
    messagesPerTick = 100,
    positionUpdateInterval = 100,
    tickIntervalsMs = 100,
  ) => {
    let currentPosition = 0;
    let messagesSinceLastPositionWrite = 0;
    let continuePolling = true;
    function writePosition(position) {
      return write(subscriberId, position);
    }
    function loadPosition() {
      return readLastMessage(subscriberId).then((message) => {
        currentPosition = message ? message.seq : 0;
      });
    }
    function updateReadPosition(position) {
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
    function handleMessage(message) {
      process.exit(1); // fix types before continuing
      const handler = handlers[message.type];
      return handler ? handler(message) : Promise.resolve(true);
    }
    function processBatch(messages) {
      console.log("processing messages...", messages);
      // @ts-ignore
      return bluebird_1.default
        .each(messages, (message) =>
          handleMessage(message)
            .then(() => updateReadPosition(message.global_position))
            .catch((err) => {
              logError(message, err);
              throw err;
            }),
        )
        .then(() => messages.length);
    }
    function logError(lastMessage, error) {
      console.error(`error processing: ${streamName}`);
    }
    function tick() {
      return getNextBatchOfMessages()
        .then(processBatch)
        .catch((err) => {
          console.error("Error processing batch", err);
          stop();
        });
    }
    function poll() {
      return __awaiter(this, void 0, void 0, function* () {
        yield loadPosition();
        while (continuePolling) {
          const messagesProcessed = yield tick();
          if (messagesProcessed === 0) {
            // @ts-ignore
            yield bluebird_1.default.delay(tickIntervalsMs);
          }
        }
      });
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
exports.configureCreateSubscription = configureCreateSubscription;
