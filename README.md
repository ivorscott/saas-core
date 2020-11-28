# CQRS Event Architecture for Devpie Client

This backend uses CQRS and event sourcing. It should be used with the [devpie-client-app](https://github.com/ivorscott/devpie-client-app).

It uses [devpie-client-events](https://github.com/ivorscott/devpie-client-common-module), a shared library that spans across 
3 language specific packages: in Typescript, Python and Golang. The Typescript commands and events are the source of truth. The library 
uses [Quicktype](https://quicktype.io/) to generate code for the Golang and Python packages. Therefore, microservices can 
be developed in multiple languages while using the shared library.

### Goal

This is an experimental project for learning. If it gets serious I will discontinue this version, clone it, and use it 
as the base for a company.

Devpie Client is a business management tool for performing software development with clients. Features will include 
kanaban or agile style board management and auxiliary services like cost estimation, payments and more. 

## How Data Moves Through The System

![cqrs architecture](cqrs.png)

End users send requests to Applications. Applications write messages (commands or events) to the Messaging System in response to those requests. Services (Components) pick up those messages, perform their work, and write new messages to the Messaging System. Aggregators observe all this activity and transform these messages into View Data that Applications use to send responses to users.

## Setup 

#### Requirements

* [Docker Desktop](https://docs.docker.com/desktop/) (Kubernetes enabled)
* [Tilt](https://tilt.dev/)
* [Create Auth0 account](http://auth0.com/)
* Fork repository to automate your own Auth0 configuration.
* [Enable Auth0 Github Deployments Extension](https://auth0.com/docs/extensions/github-deployments)
    
#### Configuration
* \__infra\__ contains the kubernetes infrastructure
* \__a0\__ contains the auth0 configuration

## Usage

Run front and back ends simultaneously. For faster development don't run the [devpie-client-app](https://github.com/ivorscott/devpie-client-app) in a container/pod.

```bash
# frontend
npm start

# backend
tilt up
```

### Debugging the database

```bash
# run temporary pod to enter a database through pgcli 
k run pgcli --rm --image=devpies/pgcli -it --command -- pgcli postgres://postgres:postgres@view-db-identity-svc:5432/identity
```

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

[x] Integrate NATS Streaming

### Components

- Components are services
- Microservices are small and focused doing one thing well
- Microservices are autonomous components that encapsulate a distinct business process
- Micoservices don't share databases with other services
- Micoservices allow us to use the technology stack best suited to achieve required performance

[x] Build Identity Service (Golang)

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

[x] Identity

[ ] Projects (jsonb)

[ ] Estimations
