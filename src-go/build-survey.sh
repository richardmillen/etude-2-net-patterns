#! /usr/bin/env bash

if [ ! -d "../bin/" ]; then
    mkdir ../bin/
fi

cd apps

printf "building 'service discovery' (by survey) apps...\n"

if [[ $1 == "" || $1 == "survey-server" ]]; then
    printf "building survey-server...\n"
    cd survey-server
    go build && mv survey-server ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "survey-client" ]]; then
    printf "building survey-client...\n"
    cd survey-client
    go build && mv survey-client ../../../bin/
    cd ..
fi

printf "finished building 'service discovery' (by survey) apps.\n\n"
