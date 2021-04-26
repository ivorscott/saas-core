import { env } from "./env";
import { connectNats } from "./msg";
import { IdentityCommandListener } from "./cmd-listener";
import { Pool } from "pg";

async function main() {
  const stan = await connectNats();

  const db = new Pool({
    connectionString: env.NATS_DB_URL,
    ssl: { rejectUnauthorized: false },
  });
  
  new IdentityCommandListener(stan, db).listen();

  stan.on("close", () => {
    console.log("NATS Streaming connection closed!");
  });

  process.on("SIGINT", () => {
    console.log("SIGINT detected.");
    stan.close();
  });

  process.on("SIGTERM", () => {
    console.log("SIGTERM detected.");
    stan.close();
  });
}

main();
