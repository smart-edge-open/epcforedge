
#! /bin/sh
#########################################################
# <COPYRIGHT_TAG>
#########################################################
setup_dir=${PWD}

set -e

curl -v --cacert mec.crt -X DELETE https://mec.local:8080/userplanes/666
