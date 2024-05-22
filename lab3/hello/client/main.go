package main

import (
	"context"
	"github.com/noartem/labs/4/2/cloud/3/hello/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Ошибка при закрытии соединения: %v", err)
		}
	}(conn)

	c := proto.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("Ошибка при вызове SayHello: %v", err)
	}

	log.Printf("Результат вызова SayHello: %s", r.GetMessage())
}
