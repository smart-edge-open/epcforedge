#! /bin/sh
setup_dir=${PWD}

set -e

curl -X POST -i "Content-Type: application/json" --data @./json/100-AF_NF_SUB_POST003.json http://localhost:8080/CNCA/1.0.1/subscriptions

exit 0

