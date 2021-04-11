import nats, { Stan } from "node-nats-streaming";
import { env } from "./env";
import { v4 as uuidV4 } from "uuid";

export let natsInstance = null

export const connectNats = (): Stan => {
  const stan = nats.connect(env.CLUSTER_ID, `${env.CLIENT_ID}-${uuidV4()}`, {
    url: env.NATS_SERVER,
  });

  stan.on("connect", () => {
    console.log("Publisher connected to NATS");
  });

  stan.on('disconnect', () => {
    console.log('[NATS] disconnect');
  });

  stan.on('reconnecting', () => {
    console.log('reconnecting');
  });

  stan.on('reconnect', () => {
    console.log(`[NATS] reconnect`);
  });

  stan.on('permission_error', function(err) {
    console.error('[NATS] got a permissions error', err.message);
  });

  stan.on('connection_lost', async (error) => {
    console.log('disconnected from stan', error);
  });

  stan.on('error', (error) => {
    console.log('[Nats] error', error);
  });

  return stan
}

export function reconnectHandler() {
  natsInstance = connectNats()
  natsInstance.on('close', () => {
    console.log("NATS Streaming connection closed!");
    setTimeout(() => reconnectHandler(), 5000);
    natsInstance = null
  })

  process.on("SIGINT", () => {
    console.log("SIGINT detected.");
    natsInstance.close();
  });

  process.on("SIGTERM", () => {
    console.log("SIGTERM detected.");
    natsInstance.close();
  });
}
