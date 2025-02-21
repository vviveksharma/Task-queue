package main

import (
	"context"
	"log"
	pb "task-queue/grpc-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(":30001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)

	resp, err := client.HealthCheck(context.Background(), &pb.NoParams{})
	if err != nil {
		log.Println("error while calling the health Check api: ", err)
	}
	log.Println(resp)
}

