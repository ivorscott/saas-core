# Common Module for Event-driven Architecture

Designed for Command Query Responsibility Segregation (CQRS) and event sourcing.

![cqrs architecture](cqrs.png)

## Overview

This package serves as a shared library for all msg interfaces, commands and events in the system. It ensures consistency and correctness while implementing the event data model across Applications, Microservices, Aggregators and various programming languages.

## How it works

Devpie Client's event data model is exported as command and event enums, and msg interfaces, to enable easy lookup of the available identifiers in the system.

```typescript
import { Commands, Events } from "@devpie/client-events";

console.log(Commands.AddUser);
// AddUser

console.log(Events.UserAdded);
// UserAdded
```

Existing interfaces allow us to type check the msg body being sent, ensuring correctness of implementation.

```typescript
export interface UserAddedEvent {
  id: string;
  subject: Events.UserAdded;
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
```

## Messaging

Messaging systems allow Microservices to exchange messages without coupling them together. Some Microservices emit messages, while others listen to the messages they subscribe to.

A msg is a generic term for data that could be either a command or an event. Commands are messages that trigger something to happen (in the future). Events are messages that notify listeners about something that has happened (in the past). Publishers send commands or events without knowing the consumers that may be listening.

### Language Support

This package is written in TypeScript but converted to additional language targets. Each supported language has its own package.

Supported languages include:

- TypeScript [See package](https://www.npmjs.com/package/@devpie/client-events)
- Golang [See package](https://github.com/ivorscott/devpie-client-events/tree/main/go)
- Python [See package](https://pypi.org/project/devpie-client-events/)

### Development

Modify `src/events.ts`, the source of truth, then re-build to update all packages.

```
npm run build
```

### Release

Here's the steps to perform a manual release for Typescript, Python and Go packages (needs to be automated). Publishing Go modules relies on git tags. https://blog.golang.org/publishing-go-modules

```bash
# 1. npm run build
# 2. update ./py/setup.py version
# 3. update ./package.json version
# 4. commit changes to git

#5. build and upload python package to PYPI
cd py
python3 setup.py sdist bdist_wheel
python3 -m twine upload --skip-existing dist/*

# 6. create a new tag for release
git tag v0.0.1

# 7. push new tag
git push origin v0.0.1

# 8. push changes to remote repository
git push origin main

# 9. publish npm module
npm publish
```
