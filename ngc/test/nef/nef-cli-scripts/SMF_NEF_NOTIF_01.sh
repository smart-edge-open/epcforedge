#! /bin/sh

setup_dir=${PWD}

set -e

# curl -X POST -i "Content-Type: application/json" --data @./json/SMF_NEF_POST_01.json http://localhost:8091/3gpp-traffic-influence/v1/AF_01/subscriptions

curl -X POST -i "Content-Type: application/json" --data @./json/SMF_NEF_NOTIF_01.json http://localhost:8080/3gpp-traffic-influence/v1/notification/upf

exit 0

