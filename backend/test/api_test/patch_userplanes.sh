
#! /bin/sh
#########################################################
# <COPYRIGHT_TAG>
#########################################################
setup_dir=${PWD}

set -e

curl --cacert epc.crt -X PATCH -H "Content-Type: application/json" --data @patch_userplanes.json https://epc.oam:8080/userplanes/2 | json_reformat

exit 0
