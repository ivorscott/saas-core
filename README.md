# Devpie Client Backend

![frontend preview](docs/images/demo.png)

## Goal

This is an experimental project for learning.

Devpie Client is a business management tool for performing software development with clients. Features will include
kanaban or agile style board management and auxiliary services like cost estimation, payments and more.

- [See Figma design](https://www.figma.com/file/M0FVvRZWGUPWgJlby4UPjm/Devpie-Client?node-id=237%3A16)
- [See Frontend repository](https://github.com/devpies/devpie-client-app)
- [See Events repository](https://github.com/devpies/devpie-client-events)
- [See Infrastructure repository](https://github.com/devpies/devpie-client-infra)

### Setup

#### Requirements

- [Docker and Kubernetes](https://docs.docker.com/desktop/)
- [Tilt](https://tilt.dev/)

- AWS
- Sengrid
- Freshbooks
- [Auth0](http://auth0.com/) with [Github Deployments](https://auth0.com/docs/extensions/github-deployments) enabled

#### Configuration

`manifests/secrets.yaml` is required. Rename `manifests/secrets.sample.yaml` and provide base64 encoded credentials for all secrets.

`core/{service}/.env` is required for integration tests.  Rename `core/{service}/.env.sample` and provide plaintext secrets.

## Developement

Run front and back ends simultaneously. For faster development we don't run the [devpie-client-app](https://github.com/ivorscott/devpie-client-app) 
in a container/pod. However, we use [tilt.dev](https://tilt.dev) to manage kubernetes development in the backend.

```bash
# devpie-client-app
npm run dev

# devpie-client-core
make up
```

Build, Test, Linting and Formatting commands exist for each service. See `core/{service}/Makefile`

```bash
make test # uses exported env vars
make unit # short tests
make fmt
make lint
make vet
make build
```

### Debugging local databases

If you want autocompletion in the terminal, use `pgcli`:

```
kubectl run pgcli --rm -i -t --env=DB_URL="postgresql://postgres:postgres@nats-db-store-svc:5432/postgres" --image devpies/pgcli
```

If you prefer the old fashion way, use `kubectl exec`:

```
kubectl exec -it <pod> -- psql -h localhost -U postgres postgres
```

If you prefer having a UI to debug postgres, you use `PgAdmin`:

```bash
kubectl run pgadmin --env="PGADMIN_DEFAULT_EMAIL=test@example.com" --env="PGADMIN_DEFAULT_PASSWORD=SuperSecret" --image dpage/pgadmin4
kubectl port-forward pod/pgadmin 8888:80
```

### Migrations

Migrations exist under the following paths:

- `./nats/migrations`
- `./core/projects/schema/migrations`
- `./core/users/schema/migrations`

#### Migration Flow

1. create a `migration`
2. add sql for `up` and `down` migration files
3. `tag` an image containing the latest migrations
4. `push` image to registry

<details>
<summary>View example</summary>
<br>

```bash
cd core/projects/schema

migrate create -ext sql -dir migrations -seq create_table

docker build -t devpies/mic-db-projects-migration:v000001 ./migrations

docker push devpies/mic-db-projects-migration:v000001
```

</details>

Then apply the latest migration with `initContainers`

_Using [init containers](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/) for migrations requires having a running database beforehand. Containers in a pod will start after init containers have executed._

<details>
<summary>View example</summary>
<br>

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mic-projects-depl
spec:
  selector:
    matchLabels:
      app: mic-projects
  template:
    metadata:
      labels:
        app: mic-projects
    spec:
      containers:
        - image: devpies/mic-projects:325b1c2
          name: mic-projects
          resources:
            requests:
              cpu: "100m"
              memory: "100Mi"
            limits:
              cpu: "250m"
              memory: "250Mi"
          env:
            - name: API_WEB_PORT
              value: ":4000"
            - name: API_WEB_CORS_ORIGINS
              value: "https://localhost:3000, https://devpie.local"
            - name: API_WEB_AUTH_DOMAIN
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: auth0-domain
            - name: API_WEB_AUTH_AUDIENCE
              valueFrom:
                secretKeyRef:
                  name: secrets
                  key: auth0-audience
            - name: API_DB_USER
              value: postgres
            - name: API_DB_NAME
              value: postgres
            - name: API_DB_PASSWORD
              value: postgres
            - name: API_DB_HOST
              value: mic-db-projects-svc
            - name: API_DB_DISABLE_TLS
              value: "true"
            - name: API_NATS_URL
              value: "nats://nats-svc:4222"
            - name: API_NATS_CLIENT_ID
              value: "mic-projects"
            - name: API_NATS_CLUSTER_ID
              value: "devpie-client"
      initContainers:
        - name: schema-migration
          image: devpies/mic-db-projects-migration:v000016
          env:
            - name: DB_URL
              value: postgresql://postgres:postgres@mic-db-projects-svc:5432/postgres?sslmode=disable
          command: ["migrate"]
          args:
            ["-path", "/migrations", "-verbose", "-database", "$(DB_URL)", "up"]
```

</details>

Learn more about migrate cli [here](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md).


## Architecture

This backend uses CQRS and event sourcing sparingly.
CQRS is not an architecture. You don't use CQRS everywhere.

This backend should be used with the [devpie-client-app](https://github.com/ivorscott/devpie-client-app).
It uses [devpie-client-events](https://github.com/ivorscott/devpie-client-common-module) as a shared library to generate
message interfaces across multiple programming languages, but the Typescript definitions in the events repository are the source of truth.

## How Data Moves Through System Parts

Two architectural models are adopted: _a traditional microservices model_ and
_an event sourcing model_ driven by CQRS.
[CQRS allows you to scale your writes and reads separately](https://medium.com/@hugo.oliveira.rocha/what-they-dont-tell-you-about-event-sourcing-6afc23c69e9a). For example, the `accounting` integration will make use CQRS to write data to Freshbooks and read data from a cache. This introduces eventual consistency and requires the frontend's support in handling eventual consistent data intelligently.

In the traditional microservices model every microservice has its own database. Within
the event sourcing model the authoritative source of truth is stored in a single message store (NATS).

In both models, messages are persisted in a message store. In the traditional microservices model,
the message store serves to promote a fault tolerant system. Microservices can have temporary downtime and return without
the loss of messages. In the event sourcing model, the current state of an entity is achieved through folding or replaying the events and running projections on an event stream. Commands exist in their own stream and we do not apply projections on them.

![two models](docs/images/arch.png)

### The Event Sourcing Model

Devpie Client will use event sourcing and CQRS sparingly. Event sourcing is useful when we want the primary source of truth to be a stream of events. This allows you to rebuild the system to match any earlier point in time. Additional benefits are increased auditing and performance but the cost is extra complexity.

Under this model, end users send requests to Applications. Applications write messages (commands or events) to the Messaging System in response to those requests. Message Handlers pick up those messages, perform their work, and write new messages to the Messaging System. Aggregators observe all this activity and transform these messages into View Data that Applications use at a later time (eventual consistency) to send responses to users.

![two models](docs/images/cqrs.png)

<details>
<summary>Read more</summary>

### Definitions

#### Applications

- Applications are not microservices.
- An Application is a feature with its own endpoints that accepts user interaction.
- Applications provide immediate responses to user input.

#### Messaging System

- A stateful msg broker plays a central role in entire architecture.
- All state transitions will be stored by NATS Streaming in streams of messages. These state transitions become the authoritative state used to make decisions.
- NATS Streaming is a durable state store as well as a transport mechanism.

#### Message Handlers

- Message Handlers are not microservices.
- Message Handlers don't have their own dedicated database.
- Message Handlers derive authoritative state from a message store using projections.
- Message Handlers are small and focused doing one thing well.

#### Aggregators

- Aggregators poll databases and watch for new events.
- Aggregators use new event data to update the View Data .
- Aggregators aggregate state transitions into View Data for Applications to retrieve at a later time.

#### View Data

- View Data are read-only models derived from state transitions.
- View Data are eventually consistent.
- View Data are not for making decisions.
- View Data are not authoritative state, but derived from authoritative state.
- View Data can be stored in any format or database that makes sense for the Application.
</details>

## Contribute

Reach out on twitter or email me if you have any questions about contributing. 

- [ivorsco77](https://twitter.com/ivorsco77)
- ivor@devpie.io
