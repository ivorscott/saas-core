# 1. FROM sets the base image to use for subsequent instructions
# Use the golang alpine image as the base stage of a multi-stage routine
FROM golang:1.14-alpine as base

# 2. ENV CGO_ENABLED=0 diables CGO which is required for production builds
ENV CGO_ENABLED=0 

# 3. WORKDIR sets the working directory for any subsequent COPY, CMD, or RUN instructions
# Set the working directory to /api
WORKDIR /

# 4. Extend the base stage and create a new stage named dev
FROM base as dev

# 5. COPY copies files or folders from source to the destination path in the image's filesystem
# Copy the go.mod and go.sum files to /api in the image's filesystem
COPY . .

# 6. Install go module dependencies in the image's filesystem
RUN go mod download

# 7. ENV sets an environment variable
# Create GOPATH and PATH environment variables
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# 8. Print go environment for debugging purposes
RUN go env && go build -o main .

CMD ["./main"]


