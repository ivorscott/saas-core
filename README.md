# Core Event-driven Architecture for Devpie Client

## Goal

This is an experimental project for learning.

Devpie Client is a business management tool for performing software development with clients. Features will include 
kanaban or agile style board management and auxiliary services like cost estimation, payments and more. 


### Setup 

#### Requirements
* [Docker Desktop](https://docs.docker.com/desktop/) (Kubernetes enabled)
* [Pgcli](https://www.pgcli.com/install)
* [Migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
* [Tilt](https://tilt.dev/)

#### Resources Being Used (no setup required)

* [2 Free 20MB Managed Database Services](elephantsql.com)
* [An Auth0 Account](http://auth0.com/)
* [Auth0 Github Deployments Extension](https://auth0.com/docs/extensions/github-deployments)
    
#### Configuration
* \__infra\__ contains the kubernetes infrastructure
* \__auth0\__ contains the auth0 configuration

Required secrets files:

1. `__infra__/secrets.yaml`
2. `.env`

Available samples:
1. `__infra__/secrets.sample.yaml`
2. `.env.sample`

Copy the secret environment variables to your `.bashrc` or `.zshrc` file. These environment variables contain remote 
database connection strings you can use to connect with pgcli for debugging.
 
If you wish, you can create aliases for routine tasks as well.

```bash
# Use inside your .bashrc or .zshrc file

export MSG_NATS=postgres://username:password@remote-db-host:5432/dbname
export MIC_DB_IDENTITY=postgres://username:password@remote-db-host:5432/dbname
export VIEW_DB_IDENTITY=postgres://username:password@remote-db-host:5432/dbname

alias k=kubectl
```

## Architecture 

This backend uses CQRS and event sourcing sparingly. 
CQRS is not an architecture. You don't use CQRS everywhere.
 
This backend should be used with the [devpie-client-app](https://github.com/ivorscott/devpie-client-app).
It uses [devpie-client-events](https://github.com/ivorscott/devpie-client-common-module) as a shared library to generate 
message interfaces across multiple programming languages, but the Typescript definitions in the events repository are the source of truth.

## How Data Moves Through CQRS System Parts

[CQRS allows you to scale your writes and reads separately](https://medium.com/@hugo.oliveira.rocha/what-they-dont-tell-you-about-event-sourcing-6afc23c69e9a). For example, the `identity` feature makes strict use CQRS to write data in one shape and read it in one or more other shapes. This introduces eventual consistency and requires the frontend's support in handling eventual consistent data intelligently. 

Devpie Client will need to display all kinds of screens to its users. End users send requests to Applications. Applications write messages (commands or events) to the Messaging System in response to those requests. Microservices pick up those messages, perform their work, and write new messages to the Messaging System. Aggregators observe all this activity and transform these messages into View Data that Applications use to send responses to users.


![cqrs pattern](cqrs.png)


### Concepts

Let's review the concepts involved in the above diagram.

<details>
<summary>Read more</summary>
<br>

#### Applications

- Applications are not microservices.
- An Application is a feature with its own endpoints that accepts user interaction.
- Applications provide immediate responses to user input.

#### Messaging System

- A stateful msg broker plays a central role in entire architecture.
- All state transitions will be stored by NATS Streaming in streams of messages. These state transitions become the authoritative state used to make decisions.
- NATS Streaming is a durable state store as well as a transport mechanism.

#### Microservices

- Microservices are autonomous units of functionality that model a single business concern.
- Microservices are small and focused doing one thing well.
- Micoservices don't share databases with other Microservices.
- Micoservices allow us to use the technology stack best suited to achieve required performance.

#### Aggregators

- Aggregators aggregate state transitions into View Data that Applications use to render a template (SSR) or enrich the client.

#### View Data

- View Data are read-only models derived from state transitions.
- View Data are eventually consistent
- View Data are not for making decisions
- View Data are not authoritative state, but derived from authoritative state.
- View Data can be stored in any format or database that makes sense for the Application
</details>

## How Data Moves Through Microservices

Most features will be plain Microservices. Microservices send and receive messages through the NATS Streaming library.

## Developement

Run front and back ends simultaneously. For faster development don't run the [devpie-client-app](https://github.com/ivorscott/devpie-client-app) in a container/pod.

```bash
# devpie-client-app
npm start

# devpie-client-cqrs-core
tilt up
```

### Testing

Navigate to the feature folder to run tests.
```bash
cd identity/application
npm run tests
```

### Debugging
 
#### Inspecting Managed Databases
Provide `pgcli` a remote connection string.
```bash
pgcli $MIC_DB_IDENTITY 

# opts: [ $MSG_NATS | $MIC_DB_IDENTITY | $VIEW_DB_IDENTITY ... ]
```

#### Using PgAdmin 
If you prefer a UI to debug postgres you may use pgadmin. Run a pod instance and then apply port fowarding. To access pgadmin go to `localhost:8888` and enter the credentials below.
```bash
kubectl run pgadmin --env="PGADMIN_DEFAULT_EMAIL=test@example.com" --env="PGADMIN_DEFAULT_PASSWORD=SuperSecret" --image dpage/pgadmin4 
kubectl port-forward pod/pgadmin 8888:80 
```
### Migrations
Microservices and Aggregators should have remote [database services](elephantsql.com).

Migrations exist under the following paths:

- `<feature>/microservice/migrations`
- `<feature>/aggregator/migrations` 

#### Migration Flow
1. move to a feature's `microservice` or `aggregator`
2. create a `migration`
3. add sql for `up` and `down` migration files
4. `tag` an image containing the latest migrations
5. `push` image to registry

<details>
<summary>View example</summary>
<br>

```bash
cd identity/microservice

migrate create -ext sql -dir migrations -seq create_table 

docker build -t devpies/mic-db-identity-migration:v000001 ./migrations

docker push devpies/mic-db-identity-migration:v000001  
```
</details>

Then apply the latest migration with `initContainers`
<details>
<summary>View example</summary>
<br>

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
</details>

Learn more about migrate cli [here](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md). 
