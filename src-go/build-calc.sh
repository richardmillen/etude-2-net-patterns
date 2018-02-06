#! /usr/bin/env bash

if [ ! -d "../bin/" ]; then
    mkdir ../bin/
fi

cd cmd

printf "building 'client/server' calculator apps...\n"

if [[ $1 == "" || $1 == "calc-client" ]]; then
    printf "building calc-client...\n"
    cd calc-client
    go build && mv calc-client ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "calc-server" ]]; then
    printf "building calc-server...\n"
    cd calc-server
    go build && mv calc-server ../../../bin/
    cd ..
fi

printf "finished building 'client/server' calculator apps.\n\n"
