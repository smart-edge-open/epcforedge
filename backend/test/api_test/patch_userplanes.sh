
#! /bin/sh
#########################################################
# <COPYRIGHT_TAG>
#########################################################
setup_dir=${PWD}

set -e

curl --cacert mec.crt -X PATCH -H "Content-Type: application/json" --data @patch_userplanes.json https://mec.local:8080/userplanes/2 | json_reformat

exit 0
