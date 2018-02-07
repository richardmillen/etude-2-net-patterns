#! /usr/bin/env bash

printf "building all apps...\n\n"

cd src-go/demos

# TODO: run demo build scripts.

printf "finished building all demo apps.\n\n"

cd ../..
cd src-go/examples

# TODO: run example build scripts.

cd ../..

printf "finished building all examples.\n\n"
