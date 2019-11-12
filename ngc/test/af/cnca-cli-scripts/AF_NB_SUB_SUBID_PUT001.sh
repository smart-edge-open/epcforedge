#! /bin/sh
setup_dir=${PWD}

set -e

curl -X PUT -i "Content-Type: application/json" --data @./json/300_AF_NF_SUB_SUBID_PUT001.json http://localhost:8080/CNCA/1.0.1/subscriptions/{subscriptionId}

exit 0

