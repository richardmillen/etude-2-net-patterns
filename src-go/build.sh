#! /usr/bin/env bash

if [ ! -d "../bin/" ]; then
    mkdir ../bin/
fi

cd apps

if [[ $1 == "" || $1 == "uuid-check" ]]; then
    printf "building 'uuid checker' (uuid-check)...\n"
    cd uuid-check
    go build && mv uuid-check ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "survey-server" ]]; then
    printf "building 'service discovery' server (survey-server)...\n"
    cd survey-server
    go build && mv survey-server ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "survey-client" ]]; then
    printf "building 'service discovery' client (survey-client)...\n"
    cd survey-client
    go build && mv survey-client ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "word-pub" ]]; then
    printf "building 'random word' publisher (word-pub)...\n"
    cd word-pub
    go build && mv word-pub ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "word-sub" ]]; then
    printf "building 'random word' subscriber (word-sub)...\n"
    cd word-sub
    go build && mv word-sub ../../../bin/
    cd ..
fi

printf "finished!\n"