package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
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

func setKafka(host, port string) {
	conn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatal("failed to get controller:", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatal("failed to dial controller:", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{{Topic: "test", NumPartitions: 3, ReplicationFactor: 1}}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatal("failed to create topic:", err)
	}
}

func main() {
	httpPort := os.Getenv("HTTP")
	kafkaHost := os.Getenv("HOST")
	kafkaPort := os.Getenv("PORT")
	kafkaTopic := os.Getenv("TOPIC")
	kafkaPartition, err := strconv.Atoi(os.Getenv("PARTITION"))
	if err != nil {
		log.Fatal("failed to parse PARTITION", err)
	}

	log.Println(httpPort, kafkaHost, kafkaPort, kafkaTopic, kafkaPartition)

	setKafka(kafkaHost, kafkaPort)

	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:%s", kafkaHost, kafkaPort), kafkaTopic, kafkaPartition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer conn.Close()

	writer := &Writer{conn: conn}

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

	log.Println("starting server on port", httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))
}
