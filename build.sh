#!/bin/sh
GOOS=windows GOARCH=386 go build -o echo_server.exe
go build -o echo_server
