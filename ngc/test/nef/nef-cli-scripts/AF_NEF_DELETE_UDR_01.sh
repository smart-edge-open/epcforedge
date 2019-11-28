#! /bin/sh

setup_dir=${PWD}

set -e

curl -X DELETE -i "Content-Type: application/json"  http://localhost:8091/3gpp-traffic-influence/v1/AF_02/subscriptions/11116

exit 0