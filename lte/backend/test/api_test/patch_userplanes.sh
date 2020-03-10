#! /bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
#

setup_dir=${PWD}

set -e

curl -X PATCH -H "Content-Type: application/json" --data @patch_userplanes.json http://epc.oam:8080/userplanes/2 | json_reformat

exit 0
