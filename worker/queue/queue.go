package queue

import (
	"context"
	"fmt"
	"log"
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
	go jobs.RunSummarizer(`Topics: These are like categories or feeds to which records are published. Imagine them as named channels where different types of data flow. For example, you might have topics like "user_activity," "order_updates," or "sensor_readings."
	Partitions: Each topic is divided into one or more partitions. Partitions are the fundamental units of parallelism in Kafka. They allow for concurrent processing of data within a topic, significantly increasing throughput. Think of partitions as lanes on a highway, allowing multiple cars (data streams) to travel simultaneously.
	Producers: These are the clients that publish data (records) to Kafka topics. Producers decide which topic and partition a record should be sent to. They can also choose to receive acknowledgements from Kafka to ensure that the data has been successfully written.
	Consumers: These are the clients that subscribe to Kafka topics and consume the data. Consumers can belong to consumer groups.
	Consumer Groups: A consumer group is a set of consumers that cooperate to consume data from a topic. Each consumer within a group is assigned to a different partition, ensuring that each record is processed by only one consumer in the group. This allows for distributed consumption and scalability. If a consumer in a group fails, another consumer in the same group takes over its assigned partitions.
	Brokers: These are the servers that make up the Kafka cluster. They store the data and handle requests from producers and consumers. A Kafka cluster typically consists of multiple brokers for fault tolerance and scalability.
	ZooKeeper: Kafka relies on ZooKeeper for cluster management, broker coordination, and configuration management. ZooKeeper helps Kafka maintain a consistent view of the cluster and ensures that consumers can discover the brokers responsible for the partitions they need to consume from.
	How Kafka Works:
	
	Producers publish data: Producers send records to a specific topic. Kafka appends these records to the appropriate partition(s) of that topic. The choice of partition can be based on a key associated with the record (for even distribution) or using other strategies.
	Data is stored in partitions: Records within a partition are ordered and immutable. Once a record is written, it cannot be changed. Kafka stores the data on disk, providing durability.
	Consumers subscribe to topics: Consumers specify the topic(s) they want to consume from and the consumer group they belong to.
	Kafka assigns partitions to consumers: Kafka assigns partitions to consumers within a group. Each consumer in a group is responsible for consuming from one or more partitions.
	Consumers consume data: Consumers read data from their assigned partitions. They keep track of their offset, which is the position of the last read record within a partition. This allows consumers to resume consumption from where they left off, even if they fail and restart.
	Benefits of Using Kafka:
	
	High Throughput: Kafka is designed for high throughput, capable of handling millions of messages per second. Its distributed architecture and efficient storage mechanisms enable it to handle massive data streams.
	Scalability: Kafka can be scaled horizontally by adding more brokers to the cluster. This allows it to handle increasing data volumes and consumer demand.
	Fault Tolerance: Kafka is fault-tolerant. Data is replicated across multiple brokers, ensuring that data is not lost even if some brokers fail. ZooKeeper plays a vital role in managing the cluster and handling broker failures.
	Durability: Data is persisted on disk, providing durability. Even if a broker fails, the data is not lost and can be recovered.
	Real-time Processing: Kafka enables real-time data processing. Consumers can subscribe to topics and receive data as soon as it is published, allowing for the creation of real-time applications.
	Decoupling: Kafka decouples producers and consumers. Producers don't need to know anything about the consumers, and vice versa. This simplifies application development and allows for independent scaling of producers and consumers.
	Use Cases for Kafka:
	
	Kafka's capabilities make it suitable for a wide range of use cases:
	
	Real-time Data Streaming: Processing and analyzing data streams in real-time, such as website activity, sensor data, financial transactions, and social media feeds.
	Log Aggregation: Collecting logs from multiple servers and applications and processing them for analysis and monitoring.
	Stream Processing: Building stream processing applications that perform transformations, aggregations, and other operations on data streams in real-time.
	Messaging: While not its primary purpose, Kafka can also be used as a message queue, although it's optimized for high-throughput streaming rather than traditional message queuing patterns.
	Microservices Communication: Enabling communication between microservices by acting as a message bus.
	`)
	fmt.Printf("Task ID %s processed successfully!\n", taskID)
}
