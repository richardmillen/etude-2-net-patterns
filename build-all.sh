#! /usr/bin/env bash

printf "building all apps...\n\n"

cd src-go

./build-calc.sh
./build-log.sh
./build-survey.sh
./build-word.sh

cd ..

printf "finished building all apps.\n\n"
