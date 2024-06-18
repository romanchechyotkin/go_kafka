package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	reader *kafka.Reader
}

func (r *Reader) Consume(ctx context.Context) {
	for {
		message, err := r.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				log.Println("context canceled")
				return
			}

			log.Printf("Error reading message: %s\n", err)
			continue
		}

		log.Printf("message at offset %d: %s = %s\n", message.Offset, string(message.Key), string(message.Value))
	}
}

func (r *Reader) Close() error {
	return r.reader.Close()
}

func main() {
	kafkaHost := os.Getenv("HOST")
	kafkaPort := os.Getenv("PORT")
	kafkaTopic := os.Getenv("TOPIC")
	kafkaPartition, err := strconv.Atoi(os.Getenv("PARTITION"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(kafkaHost, kafkaPort, kafkaTopic, kafkaPartition)

	reader := kafka.NewReader(kafka.ReaderConfig{
		StartOffset: 0,
		Brokers:     []string{fmt.Sprintf("%s:%s", kafkaHost, kafkaPort)},
		Topic:       kafkaTopic,
		Partition:   kafkaPartition,
	})
	reader.SetOffset(kafka.FirstOffset)
	r := &Reader{reader: reader}

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Println("shutting down on signal", sig)
		reader.Close()
		cancel()
	}()

	r.Consume(ctx)
}
