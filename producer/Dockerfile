FROM golang:1.22-alpine AS builder

ENV HTTP=""
ENV HOST=""
ENV PORT=""
ENV TOPIC=""
ENV PARTITION=""

RUN apk --no-cache add git

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -v -o /app/bin/producer producer/main.go

FROM alpine AS runner

WORKDIR /

COPY --from=builder /app/bin/producer .

RUN chmod +x /producer
CMD ["./producer"]