
up:
	kubectl apply -f ./databases
	tilt up
.PHONY: up

down: 
	kubectl delete -f ./databases
	tilt down
.PHONY: down

db: 
	kubectl apply -f ./databases
.PHONY: db

dbd:
	kubectl delete -f ./databases
.PHONY: dbd