
#! /bin/sh
#########################################################
# <COPYRIGHT_TAG>
#########################################################
setup_dir=${PWD}

set -e

curl --cacert mec.crt -X POST -H "Content-Type: application/json" --data @post_userplanes.json https://mec.local:8080/userplanes | json_reformat

exit 0
