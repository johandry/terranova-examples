# Builder image
FROM golang AS builder

ARG APP_NAME=simple

ENV GO111MODULE=on

WORKDIR /workspace/${APP_NAME}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /${APP_NAME} 

# Application image
FROM alpine:3.9 AS application

COPY --from=builder /simple /app/

ENTRYPOINT [ "/app/simple" ]
