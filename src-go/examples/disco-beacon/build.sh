#! /usr/bin/env bash

if [ ! -d "../../../bin/" ]; then
    mkdir ../../../bin/
fi

cd cmd

printf "building 'foo'...\n"

if [[ $1 == "" || $1 == "..." ]]; then
    printf "building 'foo'...\n"
    cd foo
    go build && mv foo ../../../../../bin/
    cd ..
fi

printf "finished building 'foo'.\n\n"
