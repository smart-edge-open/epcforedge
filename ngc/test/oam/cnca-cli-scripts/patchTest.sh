#! /bin/sh

setup_dir=${PWD}

set -e

#curl -X PATCH -H "Content-Type: application/json" --data @./json/POST001.json http://localhost:8080/afRegisters
curl -X PATCH -i "Content-Type: application/json" --data @./json/POST001.json http://localhost:8080/afRegisters/5

exit 0
