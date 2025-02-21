package queue

import (
	"context"
	"fmt"
	"log"
	"task-queue/worker/models"
	"task-queue/worker/jobs"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func Queue(dbConn *gorm.DB) {
	brokerAddress := "localhost:9092"
	topic := "task-queue"
	groupID := "task-consumer-group"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.FirstOffset,
	})

	fmt.Println("Kafka consumer started... Listening for messages...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("Error reading message:", err)
			break
		}

		fmt.Printf("Received Task ID: %s\n", string(msg.Key))
		processTask(dbConn, string(msg.Key))
		if err := reader.CommitMessages(context.Background(), msg); err != nil {
			log.Println("Failed to commit message:", err)
		}
	}
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
	go func ()  {
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
