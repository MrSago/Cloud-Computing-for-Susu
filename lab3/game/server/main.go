package main

import (
	"github.com/noartem/labs/4/2/cloud/3/game/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	gameServer := NewGameServer()
	go gameServer.StartEventLoop()

	s := grpc.NewServer()
	proto.RegisterGameServer(s, gameServer)

	log.Printf("Сервер запущен на порту %v. Ожидание подключений...", listener.Addr())

	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Ошибка при обработке запросов: %v", err)
	}
}
