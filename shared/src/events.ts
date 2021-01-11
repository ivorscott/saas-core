// System messages are either commands or events

// Events are things that have happened. (in the past)

// All messages have an id, a subject, metadata and data

export type EventTypes = Events

export interface Message {
  id: string;
  type: EventTypes;
  metadata: Metadata;
  data: any;
}

export interface Metadata {
  traceId: string;
  userId: string;
}

export enum Categories {
  Identity = "identity",
  Estimation = "estimation",
  Projects = "projects",
}

// Entity streams
// `{Category.Identity}.123`
// `{Category.Identity}.133`
// `{Category.Identity}.124`

// Category streams
// `{Category.Identity}`

// Command streams
// `{Category.Identity}.command.123`

// Command Category streams
// `{Category.Identity}.command`

export enum Commands {
  AddUser = "AddUser",
  ModifyUser = "ModifyUser",
  LockAccount = "LockAccount",
  UnlockAccount = "UnlockAccount",
  CloseAccount = "CloseAccount",
  // AddTask = 'AddTask',
  // ModifyTask = 'ModifyTask',
  // DestroyTask = 'DestroyTask',
  // MoveTask = 'MoveTask',
  // AddColumn = 'AddColumn',
  // ModifyColumn = 'ModifyColumn',
  // DestroyColumn = 'DestroyColumn',
  // AddProject = 'AddProject',
  // ModifyProject = 'ModifyProject',
  // DestroyProject = 'DestroyProject',
}

export interface AddUserCommand {
  id: string;
  type: Commands.AddUser;
  metadata: Metadata;
  data: {
    id: string;
    auth0Id: string;
    email: string;
    emailVerified: boolean;
    firstName: string;
    lastName: string;
    picture: string;
    locale: string;
  };
}

export interface ModifyUserCommand {
  id: string;
  type: Commands.ModifyUser;
  metadata: Metadata;
  data: {
    firstName: string;
    lastName: string;
    picture: string;
    locale: string;
  };
}

export enum Events {
  UserAdded = "UserAdded",
  UserModified = "UserModified",
  AccountClosed = "AccountClosed",
  AccountLocked = "AccountLocked",
  AccountUnlocked = "AccountUnlocked",
  // TaskAdded = 'TaskAdded',
  // TaskModified = 'TaskModified',
  // TaskDestroyed = 'TaskDestroyed',
  // TaskMoved = 'TaskMoved',
  // ColumnAdded = 'ColumnAdded',
  // ColumnModified= 'ColumnModified',
  // ColumnDestroyed = 'ColumnDestroyed',
  // ProjectAdded = 'ProjectAdded',
  // ProjectModified= 'ProjectModified',
  // ProjectDestroyed = 'ProjectDestroyed'
}

export interface UserAddedEvent {
  id: string;
  type: Events.UserAdded;
  metadata: Metadata;
  data: {
    id: string;
    auth0Id: string;
    email: string;
    emailVerified: boolean;
    firstName: string;
    lastName: string;
    picture: string;
    locale: string;
  };
}

export interface UserModifiedEvent {
  id: string;
  type: Events.UserModified;
  metadata: Metadata;
  data: {
    firstName: string;
    lastName: string;
    picture: string;
    locale: string;
  };
}
