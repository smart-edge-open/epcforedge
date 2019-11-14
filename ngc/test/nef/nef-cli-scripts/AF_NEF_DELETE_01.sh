#! /bin/sh

setup_dir=${PWD}

set -e

curl -X DELETE -i "Content-Type: application/json"  http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions/11111

exit 0