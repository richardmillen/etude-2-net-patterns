#! /usr/bin/env bash

if [ ! -d "../bin/" ]; then
    mkdir ../bin/
fi

cd apps

if [[ $1 == "" || $1 == "word-sub" ]]; then
    printf "building 'random word' publisher...\n"
    cd word-pub
    go build
    mv word-pub ../../../bin/
    cd ..
fi

if [[ $1 == "" || $1 == "word-sub" ]]; then
    printf "building 'random word' subscriber...\n"
    cd word-sub
    go build
    mv word-sub ../../../bin/
    cd ..
fi

printf "finished!\n"