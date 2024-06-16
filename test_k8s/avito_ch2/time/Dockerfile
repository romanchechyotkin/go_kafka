FROM golang:1.22-alpine

COPY go.mod go.sum test_k8s/avito_ch2/time/time_server.go /app/

WORKDIR /app

RUN ls -la . && \
    go install . && \
    which go_kafka

ENTRYPOINT ["/go/bin/go_kafka"]
