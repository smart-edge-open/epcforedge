#! /bin/sh
setup_dir=${PWD}

set -e

curl -X PATCH -i "Content-Type: application/json" --data @./json/400_AF_NF_SUB_SUBID_PATCH002.json http://localhost:8080/CNCA/1.0.1/subscriptions/{subscriptionId}

exit 0

