# Base image to load all the dependencies and modules
FROM golang AS base

ENV GO111MODULE=on

WORKDIR /workspace/simple

COPY ./pkg ./pkg
COPY go.mod .
COPY go.sum .

RUN  go mod download

# -----------------------------------------------
# Base image for development and test the build
# Use it with `build --target=builder-dev`
FROM base AS builder-dev

ENTRYPOINT [ "bash" ]

# -----------------------------------------------
# Image to build the application
FROM base AS builder

WORKDIR /workspace/simple
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /simple 

# -----------------------------------------------
# Application image used for development and test
# Use it with `build --target=app-dev`
FROM alpine:3.9 AS app-dev

COPY --from=builder /simple /app/

ENTRYPOINT [ "ash" ]

# -----------------------------------------------
# Application image, use it with `build --target=app`
FROM alpine:3.9 AS app

COPY --from=builder /simple /app/

ENTRYPOINT [ "/app/simple" ]