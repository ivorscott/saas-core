include .env

web:
	@CompileDaemon \
	-build="go build -o ./bin/client ./cmd/client" \
	-command="./bin/client \
	-app-backend=${API_BACKEND}" \
	-include="*.gohtml"

api:
	@CompileDaemon \
	-build="go build -o ./bin/api ./cmd/api" \
	-command="./bin/api \
	-app-frontend=${APP_FRONTEND}"
