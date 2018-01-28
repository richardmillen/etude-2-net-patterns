#! /usr/bin/env bash

if [ ! -d "../bin/" ]; then
    mkdir ../bin/
fi

cd apps

printf "building 'pub-sub' (random word publisher) apps...\n"

if [[ $1 == "" || $1 == "word-pub" ]]; then
    printf "building word-pub...\n"
    cd word-pub
    go build && mv word-pub ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "word-sub" ]]; then
    printf "building word-sub...\n"
    cd word-sub
    go build && mv word-sub ../../../bin/
    cd ..
fi

printf "finished building 'pub-sub' (random word publisher) apps.\n\n"
