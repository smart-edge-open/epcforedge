#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation



echo "Running all tests"
./cliTest.sh -m POST -i 0
./cliTest.sh -m GET -i 123457
./cliTest.sh -m PATCH -i 123457
./cliTest.sh -m GET -i 123457
./cliTest.sh -m DEL -i 123457


echo "Completed tests"
