package queue

import (
	"fmt"
	"log"
	"task-queue/worker/jobs"
	"task-queue/worker/models"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func Queue(dbConn *gorm.DB) {

	conn, err := amqp.Dial("amqp://admin:admin@localhost:5672")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		"task-queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		"task-queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	fmt.Println(" [*] Waiting for messages. To exit press CTRL+C")

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			fmt.Printf(" [x] Received: %s\n", msg.Body)
			fmt.Printf("Received Task ID: %s\n", string(msg.Body))
			processTask(dbConn, string(msg.Body))
		}
	}()

	<-forever

}

func processTask(dbConn *gorm.DB, taskID string) {
	fmt.Printf("Processing Task ID: %s...\n", taskID)
	fmt.Println("running the task please wait", dbConn)

	var taskResponse *models.DBJob

	task := dbConn.Model(&models.DBJob{}).Where("task_id = ?", taskID).First(&taskResponse)
	if task.Error != nil {
		log.Println("error while fetching the task: ", task.Error)
		return
	}
	go func() {
		resp, err := jobs.RunSummarizer(taskResponse.Data, dbConn)
		if err != nil {
			log.Println("error while running the summarizer: ", err)
			return
		}
		taskResponse.Status = "COMPLETED"
		taskResponse.Response = resp
		dbConn.Save(&taskResponse)
	}()
	fmt.Printf("Task ID %s processed successfully!\n", taskID)
}
