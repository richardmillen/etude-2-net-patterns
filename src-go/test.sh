#! /usr/bin/env bash

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

printf "finished running tests.\n\n"

