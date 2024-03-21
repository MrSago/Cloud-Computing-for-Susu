package main

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	simplejsondb "github.com/pnkj-kmr/simple-json-db"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type BulletinServer struct {
	ServeMux http.ServeMux
}

type Answer struct {
	AnswerType string
	Value      interface{}
}

func NewBulletinServer() *BulletinServer {
	server := &BulletinServer{}
	server.ServeMux.Handle("/", http.FileServer(http.Dir(".")))
	server.ServeMux.HandleFunc("/bulletin_board", server.BulletinHandler)
	return server
}

func (server *BulletinServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server.ServeMux.ServeHTTP(writer, request)
}

func (server *BulletinServer) BulletinHandler(writer http.ResponseWriter, request *http.Request) {
	opts := websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:*"},
	}

	connection, err := websocket.Accept(writer, request, &opts)
	log.Println("Connection accepted")
	if err != nil {
		panic(err)
	}

	defer connection.Close(http.StatusInternalServerError, "Server is falling")

	var item string
	var answer Answer
	for {
		err = wsjson.Read(context.Background(), connection, &item)
		if err != nil {
			panic(err)
		}

		if item == "LIST" {
			var itemsCollection simplejsondb.Collection
			itemsCollection, err = db.Collection("items")
			if err != nil {
				log.Printf("Table error: %s\n", err)
				continue
			}
			itemsString := getAllItems(itemsCollection)

			var answer Answer
			answer.AnswerType = "LIST"
			answer.Value = itemsString
			wsjson.Write(request.Context(), connection, answer)

			continue
		} else if item == "" {
			// close connection
			answer.AnswerType = "MESSAGE"
			answer.Value = "Connection closed"
			wsjson.Write(request.Context(), connection, answer)
			log.Println(answer.Value)
			break
		} else {
			var itemsCollection simplejsondb.Collection
			itemsCollection, err = db.Collection("items")
			if err != nil {
				log.Printf("Table error: %s\n", err)
				continue
			}
			addItem(itemsCollection, uuid.NewString(), item)

			answer.AnswerType = "MESSAGE"
			answer.Value = "Message added: " + item
			wsjson.Write(request.Context(), connection, answer)

			continue
		}
	}
}
