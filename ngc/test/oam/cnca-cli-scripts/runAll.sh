#!/bin/bash

echo "Running all tests"
./cliTest.sh -m POST -i 0
./cliTest.sh -m GET -i 1
./cliTest.sh -m PATCH -i 1
./cliTest.sh -m GET -i 1
./cliTest.sh -m DELDNN -i 1
./cliTest.sh -m GET -i 1
./cliTest.sh -m DEL -i 1


echo "Completed tests"
