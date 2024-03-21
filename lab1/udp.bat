@echo off
start cmd /k go run ./udp/server/main.go
timeout /t 3
start cmd /k go run ./udp/client/main.go
start cmd /k go run ./udp/client/main.go