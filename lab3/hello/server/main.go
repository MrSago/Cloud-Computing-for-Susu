package main

import (
	"context"
	"fmt"
	"github.com/noartem/labs/4/2/cloud/3/hello/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedHelloServer
}

func (s *server) SayHello(_ context.Context, in *proto.HelloRequest) (*proto.HelloReply, error) {
	log.Printf("Сообщение от клиента: %v", in.GetName())

	reply := &proto.HelloReply{
		Message: fmt.Sprintf("Hello, %s!", in.GetName()),
	}
	return reply, nil
}

func main() {
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterHelloServer(s, &server{})

	log.Printf("Сервер запущен на порту %v. Ожидание подключений...", listener.Addr())

	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Ошибка при обработке запросов: %v", err)
	}
}
