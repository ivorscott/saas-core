# Devpie Client Backend

![frontend preview](docs/images/demo.png)

## Goal

Devpie Client is a business management tool for performing software development with clients. Features will include
kanban or agile style board management and auxiliary services like cost estimation, payments and more.

- [See Figma design](https://www.figma.com/file/M0FVvRZWGUPWgJlby4UPjm/Devpie-Client?node-id=237%3A16)
- [See Frontend repository](https://github.com/devpies/devpie-client-app)
- [See Events repository](https://github.com/devpies/devpie-client-events)
- [See Infrastructure repository](https://github.com/devpies/devpie-client-infra)

### Setup

#### Requirements
- [Go](https://golang.org/doc/install) 1.13+
- Docker Desktop 3.3.3 ([Windows](https://docs.docker.com/docker-for-windows/release-notes/#docker-desktop-333), [Mac](https://docs.docker.com/docker-for-mac/release-notes/#docker-desktop-333))
- [Tilt](https://tilt.dev/) 0.20+
- [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
- [migrate](https://github.com/golang-migrate/migrate/releases)
- [shadow linter ](golang.org/x/tools/go/analysis/passes/shadow) `go get golang.org/x/tools/go/analysis/passes/shadow  `
- Secrets

#### ***Note For Windows Users***

<details>
<summary> Please Enable WSL </summary>
<br>sha

![frontend preview](docs/images/windows.png)
</details>


#### Secrets

`manifests/secrets.yaml` is required for deployments.

`core/{service}/.env` is required for end to end testing of each service.

`.gitpass` is read in the [Tiltfile](https://github.com/devpies/devpie-client-core/blob/2ddeab2eace966283f55cac58aa945a62c0c8aad/Tiltfile#L11) and passed to [Dockerfiles as build args](https://github.com/devpies/devpie-client-core/blob/2ddeab2eace966283f55cac58aa945a62c0c8aad/core/users/Dockerfile#L20). This allows services to pull private go modules but currently all repositories are public.

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
make test # short tests
make e2e # uses exported env vars
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

## Contribute

Reach out on twitter or email me if you have any questions about contributing. 

- [ivorsco77](https://twitter.com/ivorsco77)
- ivor@devpie.io
