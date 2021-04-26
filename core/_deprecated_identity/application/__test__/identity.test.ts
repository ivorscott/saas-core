import { Pool } from "pg";
import nats, { Stan } from "node-nats-streaming";
import { Request, Response } from "express";
import { AddUserPublisher } from "../publish-add-user";
import { Commands } from "@devpie/client-events";
import {
  createActions,
  createQueries,
  createIdentity,
  SQL,
  ERR,
  errors,
  sqlStatements,
} from "../identity";
import {
  formatUserRecords,
  getFakeAuth0Users,
  getFakeUsersRecords,
} from "./factory";
import mock = jest.mock;

jest.mock("../publish-add-user");
jest.mock("express", () => {
  const mockRoute = {
    get: jest.fn(),
    post: jest.fn(),
    delete: jest.fn(),
    put: jest.fn(),
  };
  const mockedRouter = {
    route: jest.fn(() => mockRoute),
  };
  return {
    Router: jest.fn(() => mockedRouter),
  };
});

jest.mock("pg", () => {
  const mockedPgClient = {
    connect: jest.fn(),
    query: jest.fn(),
    end: jest.fn(),
  };
  return { Pool: jest.fn(() => mockedPgClient) };
});

jest.mock("node-nats-streaming", () => {
  const mockedSTANClient = {
    close: jest.fn(),
    publish: jest.fn(),
    subscribe: jest.fn(),
    subscriptionOptions: jest.fn(),
  };
  const client = (mockedSTANClient as unknown) as Stan;
  return { connect: jest.fn().mockImplementation(() => client) };
});

let dbPool: Pool, natsClient: Stan;

const mockRequest = ({ auth0Id, ...rest }) => {
  return {
    user: { sub: auth0Id },
    ...rest,
  } as Request;
};

const mockResponse = () => {
  let res = {
    status: jest.fn(),
    send: jest.fn(),
    end: jest.fn(),
  };
  res.status = jest.fn().mockReturnValue(res);
  res.send = jest.fn().mockReturnValue(res);
  res.end = jest.fn().mockReturnValue(res);
  return (res as unknown) as Response;
};

beforeEach(() => {
  dbPool = new Pool({ connectionString: "DATABASE_URL" });
  natsClient = nats.connect("CLUSTER_ID", "CLIENT_ID", { url: "NATS_SERVER" });
  (AddUserPublisher as jest.Mock).mockClear();
  (dbPool.query as jest.Mock).mockClear();
});

describe("Test Actions", () => {
  test("getUser() retrieves user with camelcase keys", async () => {
    const [record] = getFakeUsersRecords();
    const [user] = formatUserRecords([record]);

    const mockQueries = {
      loadUser: jest.fn().mockResolvedValueOnce(user),
    };

    const actions = createActions(natsClient, mockQueries);
    const result = await actions.getUser(record.auth0_id);

    expect(mockQueries.loadUser).toBeCalledTimes(1);
    expect(mockQueries.loadUser).toBeCalledWith(record.auth0_id);
    expect(result).toEqual(user);
  });

  test("addUser() publishes command correctly", async () => {
    const traceId = "bbc31335-8hds-4ae9-j729-dc1efb1ccfed";

    const [record] = getFakeUsersRecords();
    const [auth0User] = getFakeAuth0Users([record]);

    const mockQueries = {
      loadUser: jest.fn(),
    };

    const mockAddUserPublisher = AddUserPublisher as jest.Mock;
    const actions = createActions(natsClient, mockQueries);

    await actions.addUser(traceId, auth0User);

    expect(mockAddUserPublisher).toHaveBeenCalledTimes(1);

    const mockAddUserPublisherInstance = mockAddUserPublisher.mock.instances[0];
    const publishCalledWith =
      mockAddUserPublisherInstance.publish.mock.calls[0][0];

    expect(publishCalledWith.type).toEqual(Commands.AddUser);
    expect(publishCalledWith.metadata).toEqual(
      expect.objectContaining({ traceId }),
    );
    expect(publishCalledWith.data).toEqual(expect.objectContaining(auth0User));
  });
});

describe("Test Queries", () => {
  test("loadUser() executes correct SQL query", async () => {
    const [record] = getFakeUsersRecords();

    (dbPool.query as jest.Mock).mockResolvedValueOnce({
      rows: [],
    });

    const queries = createQueries(dbPool);
    await queries.loadUser(record.user_id);

    expect(dbPool.query).toBeCalledTimes(1);

    const functionCalledWith = (dbPool.query as jest.Mock).mock.calls[0];
    const queryText = functionCalledWith[0];
    const values = functionCalledWith[1];

    expect(queryText).toBe(sqlStatements[SQL.GetUser]);
    expect(values).toEqual([record.user_id]);
  });
});

describe("Test Handlers", () => {
  test("findIdentity() retrieves user with 200 ok response", async () => {
    const [record] = getFakeUsersRecords();
    const [user] = formatUserRecords([record]);

    (dbPool.query as jest.Mock).mockResolvedValueOnce({
      rows: [record],
    });

    const { handlers } = createIdentity(dbPool, natsClient);
    const req = mockRequest({ auth0Id: record.auth0_id });
    const res = mockResponse();

    await handlers.findIdentity(req, res);

    expect(res.status).toHaveBeenCalledWith(200);
    expect(res.status).toHaveBeenCalledTimes(1);
    expect(res.send).toHaveBeenCalledWith(user);
  });

  test("findIdentity() returns 'user not found' with 404 response", async () => {
    const [record] = getFakeUsersRecords();
    const [user] = formatUserRecords([record]);

    (dbPool.query as jest.Mock).mockResolvedValueOnce({
      rows: [],
    });

    const { handlers } = createIdentity(dbPool, natsClient);
    const req = mockRequest({ auth0Id: record.auth0_id });
    const res = mockResponse();

    await handlers.findIdentity(req, res);

    expect(res.status).toHaveBeenCalledWith(404);
    expect(res.status).toHaveBeenCalledTimes(1);
    expect(res.send).toHaveBeenCalledWith(errors[ERR.UserNotFound]);
  });
});
