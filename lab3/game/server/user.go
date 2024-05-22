package main

import "github.com/google/uuid"

type User struct {
	Id   string
	Name string
}

func NewUser(name string) *User {
	return &User{
		Id:   uuid.New().String(),
		Name: name,
	}
}
