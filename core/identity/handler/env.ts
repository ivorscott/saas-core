enum Env {
  NODE_ENV = "NODE_ENV",
  VERSION = "VERSION",
  NATS_SERVER = "NATS_SERVER",
  NATS_DB_URL = "NATS_DB_URL",
  CLUSTER_ID = "CLUSTER_ID",
  CLIENT_ID = "CLIENT_ID",
  TIMEZONE = "TIMEZONE",
}

export interface Environment {
  [Env.NODE_ENV]: string;
  [Env.NATS_SERVER]: string;
  [Env.NATS_DB_URL]: string;
  [Env.CLUSTER_ID]: string;
  [Env.CLIENT_ID]: string;
  [Env.TIMEZONE]: string;
}

type Maybe<T> = T | undefined;

function checkValue<T>(value: Maybe<T>, name: string): T {
  if (!value) {
    throw new ReferenceError(`${name} is undefined`);
  }
  return value;
}

function importFromEnv(key: string): string {
  const value = process.env[key];
  return checkValue<string>(value, key);
}

const env: Environment = {
  [Env.NODE_ENV]: importFromEnv(Env.NODE_ENV),
  [Env.NATS_SERVER]: importFromEnv(Env.NATS_SERVER),
  [Env.NATS_DB_URL]: importFromEnv(Env.NATS_DB_URL),
  [Env.CLUSTER_ID]: importFromEnv(Env.CLUSTER_ID),
  [Env.CLIENT_ID]: importFromEnv(Env.CLIENT_ID),
  [Env.TIMEZONE]: importFromEnv(Env.TIMEZONE),
};

export { env };
