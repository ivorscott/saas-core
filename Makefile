default: up

build:
	docker build -t devpies/mic-db-users-migration:v000001 ./core/users/schema/migrations
	docker build -t devpies/mic-db-projects-migration:v000001 ./core/projects/schema/migrations
	docker build -t devpies/msg-db-nats-migration:v000003 ./nats/migrations

up: build
	kubectl config use-context docker-desktop
	kubectl apply -f ./manifests/db-nats-depl.yaml
	kubectl apply -f ./manifests/db-projects-depl.yaml
	kubectl apply -f ./manifests/db-users-depl.yaml
	tilt up
.PHONY: up

down: 
	kubectl delete -f ./manifests/db-nats-depl.yaml
	kubectl delete -f ./manifests/db-projects-depl.yaml
	kubectl delete -f ./manifests/db-users-depl.yaml
	tilt down
.PHONY: down

db: 
	kubectl apply -f ./databases
.PHONY: db

dbd:
	kubectl delete -f ./databases
.PHONY: dbd
