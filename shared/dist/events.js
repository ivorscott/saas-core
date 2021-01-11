"use strict";
// System messages are either commands or events
Object.defineProperty(exports, "__esModule", { value: true });
exports.Events = exports.Commands = exports.Categories = void 0;
var Categories;
(function (Categories) {
    Categories["Identity"] = "identity";
    Categories["Estimation"] = "estimation";
    Categories["Projects"] = "projects";
})(Categories = exports.Categories || (exports.Categories = {}));
// Entity streams
// `{Category.Identity}-123`
// Category streams
// `{Category.Identity}`
// Command streams
// `{Category.Identity}:command-123`
// Command Category streams
// `{Category.Identity}:command`
var Commands;
(function (Commands) {
    Commands["AddUser"] = "AddUser";
    Commands["ModifyUser"] = "ModifyUser";
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
})(Commands = exports.Commands || (exports.Commands = {}));
var Events;
(function (Events) {
    Events["UserAdded"] = "UserAdded";
    Events["UserModified"] = "UserModified";
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
})(Events = exports.Events || (exports.Events = {}));
