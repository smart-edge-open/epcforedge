#SPDX-License-Identifier: Apache-2.0
#Copyright © 2019 Intel Corporation

#! /bin/sh

setup_dir=${PWD}

set -e

curl -X POST -i "Content-Type: application/json" --data @./json/100_AF_NB_SUB_POST01.json http://localhost:8080/af/v1/subscriptions
curl -X POST -i "Content-Type: application/json" --data @./json/100_AF_NB_SUB_POST02.json http://localhost:8080/af/v1/subscriptions
curl -X POST -i "Content-Type: application/json" --data @./json/100_AF_NB_SUB_POST03.json http://localhost:8080/af/v1/subscriptions
curl -X POST -i "Content-Type: application/json" --data @./json/100_AF_NB_SUB_POST04.json http://localhost:8080/af/v1/subscriptions
curl -X POST -i "Content-Type: application/json" --data @./json/100_AF_NB_SUB_POST05.json http://localhost:8080/af/v1/subscriptions

exit 0

