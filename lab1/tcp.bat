@echo off
start cmd /k go run ./tcp/server/main.go
timeout /t 3
start cmd /k go run ./tcp/client/main.go
start cmd /k go run ./tcp/client/main.go