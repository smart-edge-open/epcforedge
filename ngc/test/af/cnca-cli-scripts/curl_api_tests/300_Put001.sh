#! /bin/sh
setup_dir=${PWD}

set -e

curl -X PUT -i "Content-Type: application/json" --data @./json/300_AF_NB_SUB_SUBID_PUT001.json http://localhost:8080/af/v1/subscriptions/101

exit 0

