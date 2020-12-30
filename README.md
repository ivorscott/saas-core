# CQRS Event Architecture for Devpie Client

This backend uses CQRS and event sourcing. It should be used with the [devpie-client-app](https://github.com/ivorscott/devpie-client-app).

It uses [devpie-client-events](https://github.com/ivorscott/devpie-client-common-module), a shared library that spans across 
3 language specific packages: in Typescript, Python and Golang. The Typescript commands and events are the source of truth. The library 
uses [Quicktype](https://quicktype.io/) to generate code for the Golang and Python packages. Therefore, microservices can 
be developed in multiple languages while using the shared library.

### Goal

This is an experimental project for learning.

Devpie Client is a business management tool for performing software development with clients. Features will include 
kanaban or agile style board management and auxiliary services like cost estimation, payments and more. 

## How Data Moves Through The System

![cqrs architecture](cqrs.png)

End users send requests to Applications. Applications write messages (commands or events) to the Messaging System in response to those requests. Microservices pick up those messages, perform their work, and write new messages to the Messaging System. Aggregators observe all this activity and transform these messages into View Data that Applications use to send responses to users.

## Setup 

#### Requirements
* [Docker Desktop](https://docs.docker.com/desktop/) (Kubernetes enabled)
* [Pgcli](https://www.pgcli.com/install)
* [Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
* [Free 20MB Managed Database Services](elephantsql.com)
* [Tilt](https://tilt.dev/)
* [Create Auth0 Account](http://auth0.com/)
* Fork repository to automate your own Auth0 configuration.
* [Enable Auth0 Github Deployments Extension](https://auth0.com/docs/extensions/github-deployments)
    
#### Configuration
* \__infra\__ contains the kubernetes infrastructure
* \__auth0\__ contains the auth0 configuration

Export environment variables for connecting to remote database services inside your `.bashrc` or `.zshrc` file.
 
If you wish, you can create aliases for routine tasks as well.

```bash
# inside your .bashrc or .zshrc file

export MIC_DB_IDENTITY=postgres://username:password@remote-db-host:5432/dbname
export VIEW_DB_IDENTITY=postgres://username:password@remote-db-host:5432/dbname

alias k=kubectl
```

## Usage

Run front and back ends simultaneously. For faster development don't run the [devpie-client-app](https://github.com/ivorscott/devpie-client-app) in a container/pod.

```bash
# devpie-client-app
npm start

# devpie-client-cqrs-core
tilt up
```

### Debugging remote databases outside of kubernetes
```bash
pgcli $MIC_DB_IDENTITY
```
### Debugging local databases within kubernetes
The nats streaming server uses an sql store to persist data. In development, it's using a local database volume. 
Connect to it by doing the following:
```bash
kubectl run pgcli --rm -i -t --env=DB_URL="postgresql://postgres:postgres@nats-svc:5432/postgres" --image devpies/pgcli
```
#### Using pgadmin from within the cluster.
To use pgadmin run a pod instance and use port fowarding. To access pgadmin go to localhost:8888 and enter credentials.
```bash
kubectl run pgadmin --env="PGADMIN_DEFAULT_EMAIL=test@example.com" --env="PGADMIN_DEFAULT_PASSWORD=SuperSecret" --image dpage/pgadmin4 
kubectl port-forward pod/pgadmin 8888:80 
```
### Migrations
Microservices and Aggregators should have remote [database services](elephantsql.com).

Migrations exist under the following paths:

- `<feature>/microservice/migrations`
- `<feature>/aggregator/migrations`

#### Migration flow
1. move to a feature's `microservice` or `aggregator`
2. create a `migration`
3. add sql for `up` and `down` migration files
4. `tag` an image containing the latest migrations
5. `push` image to registry

For example:

```bash
cd identity/microservice

migrate create -ext sql -dir migrations -seq create_table 

docker build -t devpies/mic-db-identity-migration:v000001 ./migrations

docker push devpies/mic-db-identity-migration:v000001  
```
Then apply the latest migration with `initContainers`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mic-identity-depl
spec:
  selector:
    matchLabels:
      app: mic-identity
  template:
    metadata:
      labels:
        app: mic-identity
    spec:
      containers:
        - image: devpies/client-mic-identity
          name: mic-identity
          env:
            - name: POSTGRES_DB
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-database-name
            - name: POSTGRES_USER
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-username
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-password
            - name: POSTGRES_HOST
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-host
# ============================================
#  Init containers are specialized containers
#  that run before app containers in a Pod.
# ============================================
      initContainers:
        - name: schema-migration
          image: devpies/mic-db-identity-migration:v000001
          env:
            - name: DB_URL
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: mic-db-identity-url
          command: ["migrate"]
          args: ["-path", "/migrations", "-verbose", "-database", "$(DB_URL)", "up"]
```
Learn more about migrate cli [here](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md). 

## Concepts

### Applications

- Applications are not microservices.
- An Application is a feature with its own endpoints that accepts user interaction.
- Applications provide immediate responses to user input.

### Messaging System

- A stateful msg broker plays a central role in entire architecture.
- All state transitions will be stored by NATS Streaming in streams of messages. These state transitions become the authoritative state used to make decisions.
- NATS Streaming is a durable state store as well as a transport mechanism.

### Microservices

- Microservices are autonomous units of functionality that model a single business concern.
- Microservices are small and focused doing one thing well.
- Micoservices don't share databases with other Microservices.
- Micoservices allow us to use the technology stack best suited to achieve required performance.

### Aggregators

- Aggregators aggregate state transitions into View Data that Applications use to render a template (SSR) or enrich the client.

### View Data

- View Data are read-only models derived from state transitions.
- View Data are eventually consistent
- View Data are not for making decisions
- View Data are not authoritative state, but derived from authoritative state.
- View Data can be stored in any format or database that makes sense for the Application
