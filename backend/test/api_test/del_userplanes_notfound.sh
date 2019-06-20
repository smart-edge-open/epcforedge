
#! /bin/sh
#########################################################
# <COPYRIGHT_TAG>
#########################################################
setup_dir=${PWD}

set -e

curl -v --cacert epc.crt -X DELETE https://epc.oam:8080/userplanes/666
