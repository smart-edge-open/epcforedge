#! /bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation


setup_dir=${PWD}
echo "$setup_dir"

set -e

#/var/tmp/curl/src/curl --http2 --cert certs/root-ca-key.pem -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_01.json https://localhost:8090/3gpp-traffic-influence/v1/AF_01/subscriptions
/var/tmp/curl/src/curl --http2 --insecure -X POST -i "Content-Type: application/json" --data @./json/AF_NEF_POST_01.json https://localhost:8090/3gpp-traffic-influence/v1/AF_01/subscriptions
exit 0

