package internal

import (
	"context"
	"fmt"
	"github.com/magomedcoder/simple-microservice/logger-service/api/pb"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type LogServer struct {
	pb.UnimplementedLogServiceServer
	Models Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
	input := req.GetLogEntry()
	logEntry := LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	if err := l.Models.LogEntry.Insert(logEntry); err != nil {
		res := &pb.LogResponse{Result: "не удалось"}
		return res, err
	}
	res := &pb.LogResponse{Result: "Записано через gRPC"}
	return res, nil
}

func (c *Config) GRPCListen() {
	port := os.Getenv("GRPC_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("ошибка при прослушивании gRPC: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLogServiceServer(s, &LogServer{Models: c.Models})
	log.Printf("gRPC сервер запущен на порту %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("ошибка при прослушивании gRPC: %v", err)
	}
}
