package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %v [ip:port]\n", os.Args[0])
		panic("Need more arguments!\n")
	}

	listener, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Server started on: %v\n", listener.Addr())

	createDb("database")
	createCollection("items")

	bulletinServer := NewBulletinServer()
	httpServer := &http.Server{
		Handler: bulletinServer,
	}

	errc := make(chan error, 1)
	go func() {
		httpServer.Serve(listener)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("не удалось обработать: %v", err)
	case sig := <-sigs:
		log.Printf("завершение: %v", sig)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return httpServer.Shutdown(ctx)
}
