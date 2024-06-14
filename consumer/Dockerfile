FROM golang:1.22-alpine AS builder

ENV HOST=""
ENV PORT=""
ENV TOPIC=""
ENV PARTITION=""

RUN apk --no-cache add git

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -v -o /app/bin/consumer consumer/main.go

FROM alpine AS runner

WORKDIR /

COPY --from=builder /app/bin/consumer .

RUN chmod +x /consumer
CMD ["./consumer"]