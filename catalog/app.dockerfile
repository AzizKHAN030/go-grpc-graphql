# Stage 1: Build Stage
FROM golang:1.23-alpine3.21 AS build

# Install necessary dependencies
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory
WORKDIR /go/src/github.com/azizkhan030/go-grpc-graphql

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download
COPY catalog ./catalog

# Build the binary for the catalog service
RUN GO111MODULE=on go build -o /go/bin/catalog ./catalog/cmd/catalog

# Stage 2: Minimal Runtime Stage
FROM alpine:3.11
WORKDIR /usr/bin

# Copy the compiled binary from the build stage
COPY --from=build /go/bin/catalog .

# Expose the application's port
EXPOSE 8080

# Set the default command to run the application
CMD ["catalog"]