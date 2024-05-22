package main

import (
	"log"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/net/websocket"
)

func sortWords(textBlock string) []string {
	words := strings.Fields(textBlock)

	wordsMap := make(map[string]bool)
	for _, word := range words {
		wordsMap[word] = true
	}
	uniqueWords := make([]string, 0, len(wordsMap))
	for word := range wordsMap {
		uniqueWords = append(uniqueWords, word)
	}
	sort.Strings(uniqueWords)
	return uniqueWords
}

func SortServer(ws *websocket.Conn) {
	buf := make([]byte, 32*1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			return
		}
		words := string(buf[:n])
		log.Println("From client: ")
		for _, word := range words {
			log.Println(word)
		}
		sortedWords := sortWords(words)
		ws.Write([]byte(strings.Join(sortedWords, " ")))
	}
}

func main() {
	http.Handle("/sort", websocket.Handler(SortServer))
	log.Println("server started")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
