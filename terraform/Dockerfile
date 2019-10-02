# Builder image
FROM golang AS builder

WORKDIR /workspace/tform

ENV     GO111MODULE=on
COPY    go.mod .
COPY    go.sum .
RUN     go mod download

COPY    . .
RUN     CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tform 

# Application image
FROM alpine:3.9 AS application

COPY --from=builder /tform /app/

ENTRYPOINT [ "/app/tform" ]
