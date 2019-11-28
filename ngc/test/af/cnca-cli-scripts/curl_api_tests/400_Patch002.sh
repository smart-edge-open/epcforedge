#! /bin/sh
setup_dir=${PWD}

set -e

curl -X PATCH -i "Content-Type: application/json" --data @./json/400_AF_NB_SUB_SUBID_PATCH002.json http://localhost:8080/af/v1/subscriptions/100

exit 0

