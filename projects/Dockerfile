FROM golang:1.15-alpine as base

ENV CGO_ENABLED=0

WORKDIR /projects

FROM base as dev

COPY go.* ./

RUN go mod download
COPY . .

RUN go env

EXPOSE 4000
RUN go build -o main ./cmd/api
CMD ["./main"]

FROM dev as build-stage

RUN go build -o main ./cmd/api

FROM aquasec/trivy:0.4.4 as trivy

RUN trivy alpine:3.13.4 && \
    echo "No image vulnerabilities" > result

FROM alpine:3.13.4 as prod

COPY --from=trivy result secure
COPY --from=build-stage /projects/main main

EXPOSE 4000

HEALTHCHECK --interval=3s --timeout=3s CMD wget --spider -q http://localhost:4000/v1/health || exit 1

CMD ["./main"]
