# Stage 1: Build Stage
FROM golang:1.23-alpine3.21 AS build

# Install necessary tools and certificates
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory inside the container
WORKDIR /go/src/github.com/azizkhan030/go-grpc-graphql

# Copy go.mod and go.sum first to leverage Docker caching for dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy vendor and account source code into the container
COPY account ./account

# Build the binary for the account service
RUN GO111MODULE=on go build -o /go/bin/account ./account/cmd/account

# Stage 2: Minimal Runtime Stage
FROM alpine:3.11
WORKDIR /usr/bin

# Copy the compiled binary from the build stage
COPY --from=build /go/bin/account .

# Expose the application's port
EXPOSE 8080

# Set the default command to run the application
CMD ["account"]