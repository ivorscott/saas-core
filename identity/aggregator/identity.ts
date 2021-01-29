import { Categories, Events, UserAddedEvent } from "@devpie/client-events";
import { Pool } from "pg";

function createHandlers(queries: any) {
  return {
    [Events.UserAdded]: async (event: any) => await queries.addUser(event),
  };
}

 function createQueries(db: Pool) {
   function addUser(event: UserAddedEvent) {
    return db.query(
      `INSERT INTO users (
        id,
        auth0_id,
        email,
        email_verified,
        first_name,
        last_name,
        picture,
        locale
      )
      VALUES ($1, $2, $3, $4, $5, $6, $7 ,$8)
      ON CONFLICT DO NOTHING`,
      [
        event.data.id,
        event.data.auth0Id,
        event.data.email,
        event.data.emailVerified,
        event.data.firstName,
        event.data.lastName,
        event.data.picture,
        event.data.locale,
      ],
    );
  }
  return {
    addUser,
  };
}

export function createAggregator(db: Pool, messageStore: any) {
  const queries = createQueries(db);
  const handlers = createHandlers(queries);

  const subscription = messageStore.createSubscription(
    Categories.Identity,
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
