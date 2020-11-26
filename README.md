# CQRS Event Architecture for Devpie Client

This project uses CQRS and event sourcing.

Event types and definitions are defined in a separate [repository](https://github.com/ivorscott/devpie-client-common-module).

The frontend application lives in its own [repository](https://github.com/ivorscott/devpie-client-app).

![cqrs architecture](cqrs.png)

### How Data Moves Through The System

End users send requests to Applications. Applications write messages (commands or events) to the Messaging System in response to those requests. Services (Components) pick up those messages, perform their work, and write new messages to the Messaging System. Aggregators observe all this activity and transform these messages into View Data that Applications use to send responses to users.

## ToDo

### Applications

- Applications are not microservices.
- An Application is a feature with its own endpoints that accepts user interaction.
- Applications provide immediate responses to user input.

[x] Build Identity App (Typescript)

[ ] Build Projects App (Typescript)

[ ] Build Estimation App (Typescript)

### Messaging System

- A stateful message broker plays a central role in entire architecture.
- All state transitions will be stored by NATS Streaming in streams of messages. These state transitions become the authoritative state we use to make decisions.
- NATS Streaming is a durable state store as well as a transport mechanism.

[ ] Integrate NATS Streaming

### Components

- Components are services
- Microservices are small and focused doing one thing well
- Microservices are autonomous components that encapsulate a distinct business process
- Micoservices don't share databases with other services
- Micoservices allow us to use the technology stack best suited to achieve required performance

[ ] Build Identity Service (Golang)

[ ] Build Projects Service (Typescript)

[ ] Build Estimation Service (Python)

### Aggregators

- Aggregators aggregate state transitions into View Data that Applications use to render a template (SSR) or enrich the client.

[ ] Identity Aggregator (Typescript)

[ ] Projects Aggregator (Typescript)

[ ] Estimations Aggregator (Typescript)

### View Data

- View Data are read-only models derived from state transtions.
- View Data are eventually consistent
- View Data are not used to make decisions
- View Data are not authoritative state, but derived from authoritative state.
- View Data can be stored in any format or database that makes sense for the Application

[ ] Identity

[ ] Projects (jsonb)

[ ] Estimations
