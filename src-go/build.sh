#! /usr/bin/env bash

if [ ! -d "../bin/" ]; then
    mkdir ../bin/
fi

cd apps

printf "building 'random word' subscriber...\n"
cd word-sub
go build
mv word-sub ../../../bin/

printf "building 'random word' publisher...\n"
cd ../word-pub
go build
mv word-pub ../../../bin/

printf "finished!\n"