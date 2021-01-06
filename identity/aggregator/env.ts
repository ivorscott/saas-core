import packageJson from "./package.json";

enum Env {
  AGG_NAME = "AGG_NAME",
  DATABASE_URL = "DATABASE_URL",
  NODE_ENV = "NODE_ENV",
  VERSION = "VERSION",
  NATS_SERVER = "NATS_SERVER",
  NATS_DB_URL = "NATS_DB_URL",
}

export interface Environment {
  [Env.AGG_NAME]: string;
  [Env.NODE_ENV]: string;
  [Env.VERSION]: string;
  [Env.NATS_SERVER]: string;
  [Env.DATABASE_URL]: string;
  [Env.NATS_DB_URL]: string;
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
  [Env.AGG_NAME]: importFromEnv(Env.AGG_NAME),
  [Env.DATABASE_URL]: importFromEnv(Env.DATABASE_URL),
  [Env.NODE_ENV]: importFromEnv(Env.NODE_ENV),
  [Env.VERSION]: packageJson.version,
  [Env.NATS_SERVER]: importFromEnv(Env.NATS_SERVER),
  [Env.NATS_DB_URL]: importFromEnv(Env.NATS_DB_URL),
};

export { env };
