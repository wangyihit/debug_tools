#!/bin/sh

s="    "
ss="${s}${s}"
echos(){
echo "${s}$*" 
}
echoss(){
echo "${ss}$*"
}
      
build_package(){
    dest=$1
    cd ../$dest
    echo "build ${dest}"
    # go tool dist list 
    echos "build for win"
    GOOS=windows GOARCH=386 go build -o ../bin/${dest}.exe
    echos "build for mac"
    GOOS=darwin GOARCH=amd64 go build -o ../bin/${des}_mac
    echos "build for pi"
    GOOS=linux GOARCH=arm go build -o ../bin/${dest}_pi
    echos "build for linux"
    go build -o ../bin/${dest}
    echo "build ${dest} complete"
    cd -
}

