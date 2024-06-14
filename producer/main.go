package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	conn *kafka.Conn
}

func (w *Writer) WriteMessages(msg ...kafka.Message) (n int, err error) {
	return w.conn.WriteMessages(msg...)
}

func main() {
	httpPort := os.Getenv("HTTP")
	kafkaHost := os.Getenv("HOST")
	kafkaPort := os.Getenv("PORT")
	kafkaTopic := os.Getenv("TOPIC")
	kafkaPartition, err := strconv.Atoi(os.Getenv("PARTITION"))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:%s", kafkaHost, kafkaPort), kafkaTopic, kafkaPartition)
	if err != nil {
		log.Fatal(err)
	}

	writer := &Writer{conn: conn}
	_, err = writer.WriteMessages()
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		bytes, _ := io.ReadAll(r.Body)

		log.Println(string(bytes))

		n, err := writer.WriteMessages(kafka.Message{
			Value: bytes,
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))

			return
		}

		log.Println(n)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))
}
