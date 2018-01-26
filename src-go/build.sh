#! /usr/bin/env bash

if [ ! -d "../bin/" ]; then
    mkdir ../bin/
fi

if [[ $1 == "" || $1 == "test" ]]; then
    printf "\nrunning package tests...\n"
    printf "\n"
    printf "\t- - - - - - - - - -\n"
    printf "\t- frames          -\n"
    printf "\t- patterns/core   -\n"
    printf "\t- patterns/disco  -\n"
    printf "\t- patterns/pubsub -\n"
    printf "\t- utils            -\n"
    printf "\t- uuid            -\n"
    printf "\t- - - - - - - - - -\n"
    printf "\n"

    cd frames
    go test -v -cover
    cd ..
    printf "\n"

    cd patterns/core
    go test -v -cover
    cd ../..
    printf "\n"

    cd patterns/disco
    go test -v -cover
    cd ../..
    printf "\n"

    cd patterns/pubsub
    go test -v -cover
    cd ../..
    printf "\n"

    cd utils
    go test -v -cover
    cd ..
    printf "\n"

    cd uuid
    go test -v -cover
    cd ..
    printf "\n"

    printf "finished running tests.\n"
    printf "\n"
fi

cd apps



if [[ $1 == "" || $1 == "log-client" ]]; then
    printf "building 'log/tracing' client (log-client)...\n"
    cd log-client
    go build && mv log-client ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "log-collector" ]]; then
    printf "building 'log/tracing' collector (log-collector)...\n"
    cd log-collector
    go build && mv log-collector ../../../bin/
    cd ..
fi



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



printf "build script completed.\n\n"
