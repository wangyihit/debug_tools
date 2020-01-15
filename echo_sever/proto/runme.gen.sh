#!/bin/sh
dest=../../pb/echo_server
mkdir -p $dest
protoc -I=.  --go_out=$dest --micro_out=$dest  \
	./echo_server.proto
