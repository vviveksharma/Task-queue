package main

import (
	"context"
	"errors"
	"log"
	"net"
	"task-queue/grpc-server/db"
	"task-queue/grpc-server/models"
	pb "task-queue/grpc-server/proto"
	"task-queue/grpc-server/queue"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

const (
	STATE_PENDING     = "PENDING"
	STATE_IN_PROGRESS = "PROGRESS"
	STATE_COMPLETED   = "COMPLETED"
)

type Server struct {
	pb.UnimplementedTaskServiceServer
	DB        *gorm.DB
	KafkaConn *kafka.Conn
}

func (s *Server) HealthCheck(ctx context.Context, in *pb.NoParams) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Message: "Connection Successfull",
	}, nil
}

func (s *Server) StartTask(ctx context.Context, req *pb.StartTaskRequest) (resp *pb.StartRequestResponse, err error) {
	if req.Id == "" {
		return nil, errors.New("error while validating the request id cannot be empty")
	}
	transaction := s.DB.Begin()
	if transaction.Error != nil {
		return nil, errors.New("error while starting the transaction" + err.Error())
	}
	defer transaction.Rollback()
	taskid := uuid.New()
	addJob := transaction.Create(&models.DBJob{
		Id:     uuid.MustParse(req.Id),
		TaskId: taskid,
		Status: STATE_PENDING,
	})

	if addJob.Error != nil {
		return nil, errors.New("error while adding the job: " + err.Error())
	}
	transaction.Commit()


	err = queue.PublishMessageIntoQueue(s.KafkaConn, taskid.String())
	if err != nil {	
		log.Println("error while publishing the message to the kafka topic: ", err)
		return nil, errors.New("error while publishing the message to the kafka: " + err.Error())
	}

	return &pb.StartRequestResponse{
		Message: taskid.String(),
		Status:  "PENDING",
	}, nil
}

func (s *Server) GetTaskStatus(ctx context.Context, req *pb.GetTaskStatusRequest) (resp *pb.GetTaskStatusResponse, err error) {
	if req.Id == "" {
		return nil, errors.New("error while validating the request id cannot be empty")
	}
	transaction := s.DB.Begin()
	if transaction.Error != nil {
		return nil, errors.New("error while starting the transaction" + err.Error())
	}
	defer transaction.Rollback()
	var statusResponse models.DBJob
	status := transaction.Find(&statusResponse, models.DBJob{
		TaskId: uuid.MustParse(req.Id),
	})
	if status.Error != nil {
		return nil, errors.New("error while finding the job status: " + err.Error())
	}
	return &pb.GetTaskStatusResponse{
		Status: statusResponse.Status,
	}, nil
}

func main() {
	var ser Server

	qconn, err := queue.Publisher()
	if err != nil {
		log.Fatal("error while starting the kafka connection: ", err)
	}

	ser.KafkaConn = qconn

	db, err := db.NewDbRequest()
	if err != nil {
		log.Fatal("error while starting the database instance: ", err)
	}
	dbConn, err := db.InitDB()
	if err != nil {
		log.Fatal("error while creating a connection to database: ", err)
	}
	ser.DB = dbConn

	conn, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal("error while starting the GRPC server: ", err)
	}
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &ser)
	log.Println("Server running at port 3000")
	s.Serve(conn)
}
