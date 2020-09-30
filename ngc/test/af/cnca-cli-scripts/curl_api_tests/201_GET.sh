#! /bin/sh
#SPDX-License-Identifier: Apache-2.0
#Copyright © 2019 Intel Corporation

setup_dir=${PWD}
echo "$setup_dir"
set -e

curl http://localhost:8080/AF/v1/subscriptions/11112

exit 0

