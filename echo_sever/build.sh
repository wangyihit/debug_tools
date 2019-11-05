#!/bin/sh
# go tool dist list 

echo "build for win"
GOOS=windows GOARCH=386 go build -o echo_server.exe
echo "build for mac"
GOOS=darwin GOARCH=amd64 go build -o echo_server_mac
echo "build for pi"
GOOS=linux GOARCH=arm go build -o echo_server_pi
echo "build for linux"
go build -o echo_server
echo "complete"

