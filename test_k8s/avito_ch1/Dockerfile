FROM golang:1.22-alpine

COPY main.go go.mod go.sum /app/

WORKDIR /app

RUN ls -la . && \
    go install . && \
    which go_kafka

ENTRYPOINT ["/go/bin/go_kafka"]
