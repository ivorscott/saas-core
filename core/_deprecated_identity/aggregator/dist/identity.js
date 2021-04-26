"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.createAggregator = void 0;
const client_events_1 = require("@devpie/client-events");
function createHandlers(queries) {
  return {
    [client_events_1.Events.UserAdded]: (event) => queries.addUser(event),
  };
}
function createQueries(db) {
  function addUser(event) {
    return db.query(
      `INSERT INTO users (
        user_id,
        auth0_id,
        email,
        email_verified,
        first_name,
        last_name,
        picture)
      VALUES ($1, $2, $3, $4, $5, $6, $7 ,$8)`,
      [
        event.data.auth0Id,
        event.data.email,
        event.data.emailVerified,
        event.data.firstName,
        event.data.lastName,
        event.data.picture,
      ],
    );
  }
  return {
    addUser,
  };
}
function createAggregator(db, messageStore) {
  console.log("creating aggregator");
  const queries = createQueries(db);
  const handlers = createHandlers(queries);
  const subscription = messageStore.createSubscription(
    client_events_1.Categories.Identity,
    handlers,
    "aggregators:user-registration",
  );
  function start() {
    subscription.start();
  }
  function stop() {
    subscription.stop();
  }
  return {
    handlers,
    queries,
    start,
    stop,
  };
}
exports.createAggregator = createAggregator;
