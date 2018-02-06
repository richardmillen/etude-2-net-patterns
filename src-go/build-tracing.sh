#! /usr/bin/env bash

if [ ! -d "../bin/" ]; then
    mkdir ../bin/
fi

cd cmd

printf "building 'logging / distributed tracing' apps...\n"

if [[ $1 == "" || $1 == "log-client" ]]; then
    printf "building log-client...\n"
    cd log-client
    go build && mv log-client ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "log-collector" ]]; then
    printf "building log-collector...\n"
    cd log-collector
    go build && mv log-collector ../../../bin/
    cd ..
fi

printf "finished building 'logging / distributed tracing' apps.\n\n"
