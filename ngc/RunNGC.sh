#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation

PID_LIST=()

# Build NGC executables. Executables of af, nef and oam created in 'dist' folder
make build

# Change Directory to 'dist' folder
cd dist

# Copy configuration files folder 'configs' to 'dist' folder
cp -r ../configs .

# Create 'certs' folder with SSL Server and Client certificate files
mkdir certs
cd certs
cp ../../scripts/genCerts.sh .
chmod +x ./genCerts.sh
./genCerts.sh -t DNS -h localhost
cd ..
# Copy 'certs' folder to /etc/
sudo cp -r ./certs /etc/

# Execute oam in backgroud mode and store pid
./oam &
PID_LIST+=($!)

# Execute af in backgroud mode and store pid
./af &
PID_LIST+=($!)

# Execute nef in backgroud mode and store pid
./nef &
PID_LIST+=($!)

function terminate()
{
    for ((idx=${#PID_LIST[@]}-1;idx>=0;idx--)); do
        kill -SIGKILL ${PID_LIST[$idx]}
    done
	cd ..
	sudo make clean
}

# Wait for the SIGINT. On SIGINT system event, terminate function will get triggered to kill af, nef and oam,
# if running

trap terminate SIGINT
wait ${PID_LIST}
