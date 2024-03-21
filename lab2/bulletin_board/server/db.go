package main

import (
	"encoding/json"
	"log"

	simplejsondb "github.com/pnkj-kmr/simple-json-db"
)

var db simplejsondb.DB

func createDb(filename string) {
	var err error
	db, err = simplejsondb.New(filename, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Database created")
}

func createCollection(collectionName string) simplejsondb.Collection {
	collection, err := db.Collection(collectionName)
	if err != nil {
		panic(err)
	}
	log.Printf("Collection %s created\n", collectionName)
	return collection
}

func addItem(collection simplejsondb.Collection, key string, data string) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Marshal error: %s\n", err)
		return err
	}

	log.Println(jsonBytes)

	err = collection.Create(key, jsonBytes)
	if err != nil {
		log.Printf("Creation error: %s\n", err)
		return err
	}

	log.Printf("Item added: %s\n", data)
	return nil
}

func getAllItems(collection simplejsondb.Collection) []string {
	jsonBytes := collection.GetAll()

	var items []string
	for _, value := range jsonBytes {
		var item string
		err := json.Unmarshal(value, &item)
		if err != nil {
			log.Printf("Unmarshal error: %s\n", err)
			return nil
		}
		items = append(items, item)
	}
	return items
}
