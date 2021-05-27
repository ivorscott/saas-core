
up:
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