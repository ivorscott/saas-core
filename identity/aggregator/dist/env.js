"use strict";
var __importDefault =
  (this && this.__importDefault) ||
  function (mod) {
    return mod && mod.__esModule ? mod : { default: mod };
  };
Object.defineProperty(exports, "__esModule", { value: true });
exports.env = void 0;
const package_json_1 = __importDefault(require("./package.json"));
var Env;
(function (Env) {
  Env["AGG_NAME"] = "AGG_NAME";
  Env["DATABASE_URL"] = "DATABASE_URL";
  Env["NODE_ENV"] = "NODE_ENV";
  Env["VERSION"] = "VERSION";
  Env["NATS_SERVER"] = "NATS_SERVER";
  Env["NATS_DB_URL"] = "NATS_DB_URL";
})(Env || (Env = {}));
function checkValue(value, name) {
  if (!value) {
    throw new ReferenceError(`${name} is undefined`);
  }
  return value;
}
function importFromEnv(key) {
  const value = process.env[key];
  return checkValue(value, key);
}
const env = {
  [Env.AGG_NAME]: importFromEnv(Env.AGG_NAME),
  [Env.DATABASE_URL]: importFromEnv(Env.DATABASE_URL),
  [Env.NODE_ENV]: importFromEnv(Env.NODE_ENV),
  [Env.VERSION]: package_json_1.default.version,
  [Env.NATS_SERVER]: importFromEnv(Env.NATS_SERVER),
  [Env.NATS_DB_URL]: importFromEnv(Env.NATS_DB_URL),
};
exports.env = env;
