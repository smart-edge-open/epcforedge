#! /bin/sh

setup_dir=${PWD}

set -e

curl -X POST -i "Content-Type: application/json" --data @./json/POST001.json http://localhost:8080/afRegisters

exit 0
