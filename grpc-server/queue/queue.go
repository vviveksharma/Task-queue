package queue

import (
	"log"

	"github.com/streadway/amqp"
)

func Publisher() (conn *amqp.Connection, err error) {

	conn, err = amqp.Dial("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	return conn, nil
}

func PublishMessageIntoQueue(conn *amqp.Connection, taskId string) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Println("error while creating channel: ", err)
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		"my_queue", // Queue name
		false,      // Durable
		false,      // Delete when unused
		false,      // Exclusive
		false,      // No-wait
		nil,        // Arguments
	)

	if err != nil {
		log.Println("error while declaring queue: ", err)
		return err
	}

	Qerr := ch.Publish(
		"",
		"task-queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(taskId),
		},
	)
	if Qerr != nil {
		log.Println("error while publishing the message on the queue: ", err)
		return err
	}
	return nil
}
