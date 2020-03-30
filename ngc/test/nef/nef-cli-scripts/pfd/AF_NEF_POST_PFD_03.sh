# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2020 Intel Corporation

#! /bin/sh

setup_dir=${PWD}

set -e

curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_PFD_03.json http://localhost:8091/3gpp-pfd-management/v1/AF_01/transactions
#curl -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_01.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions


exit 0

