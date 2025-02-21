package queue

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func Publisher() (conn *kafka.Conn, err error) {
	topic := "task-queue"
	partition := 0

	conn, err = kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("Failed to connect:", err)
		return nil, err
	}

	// defer conn.Close()

	return conn, nil
}

func PublishMessageIntoQueue(conn *kafka.Conn,taskId string) error{
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err := conn.WriteMessages(
		kafka.Message{Key: []byte(taskId)},
	)
	if err != nil {
		log.Println("Failed to write messages:", err)
		return err
	}
	return nil
}
