import nats, { Stan } from "node-nats-streaming";
import { env } from "./env";
import { v4 as uuidV4 } from "uuid";

export const connectNats = () => {
  const stan = nats.connect(env.CLUSTER_ID, `${env.CLIENT_ID}-${uuidV4()}`, {
    url: env.NATS_SERVER,
  });

  return new Promise<Stan>((resolve, reject) => {
    stan.on("connect", () => {
      console.log("Publisher connected to NATS");
      resolve(stan);
    });

    stan.on("error", (error) => {
      console.log("[Nats] error", error);
      reject(error);
    });
  });
};
