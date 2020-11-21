import packageJson from "./package.json";

enum Env {
  APP_NAME = "APP_NAME",
  DATABASE_URL = "DATABASE_URL",
  NODE_ENV = "NODE_ENV",
  PORT = "PORT",
  VERSION = "VERSION",
  AUTH0_DOMAIN = "AUTH0_DOMAIN",
  AUTH0_AUDIENCE = "AUTH0_AUDIENCE",
}

export interface Environment {
  [Env.APP_NAME]: string;
  [Env.DATABASE_URL]: string;
  [Env.NODE_ENV]: string;
  [Env.PORT]: number;
  [Env.VERSION]: string;
  [Env.AUTH0_DOMAIN]: string;
  [Env.AUTH0_AUDIENCE]: string;
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
  [Env.APP_NAME]: importFromEnv(Env.APP_NAME),
  [Env.AUTH0_AUDIENCE]: importFromEnv(Env.AUTH0_AUDIENCE),
  [Env.AUTH0_DOMAIN]: importFromEnv(Env.AUTH0_DOMAIN),
  [Env.DATABASE_URL]: importFromEnv(Env.DATABASE_URL),
  [Env.NODE_ENV]: importFromEnv(Env.NODE_ENV),
  [Env.PORT]: parseInt(importFromEnv(Env.PORT), 10),
  [Env.VERSION]: packageJson.version,
};

export { env };
