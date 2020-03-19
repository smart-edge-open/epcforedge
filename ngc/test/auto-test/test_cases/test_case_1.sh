#!/bin/bash

#SPDX-License-Identifier: Apache-2.0
#Copyright © 2019 Intel Corporation

# A basic test case which uses different function exported from the test_nef.sh
# For generic test case use see test_case_2.sh



# Importing test suite functions
source test_nef.sh

# First test case
echo "1: Testing (post json/AF_NEF_POST_02.json)"
post json/AF_NEF_POST_02.json

if [[ $status_code -ne 201 ]]; then
	echo "Test Failed"
	echo "Returned body: $body"
else
	echo "Passed"
fi

#Getting returned sub id from json
ret_id=`echo $body | jq -r '.self' | awk 'BEGIN {FS="/"} // {print $(NF)}'`

# 2nd test case
echo "2: Testing (get $ret_id)"
get $ret_id 

if [[ $status_code -ne 200 ]]; then
	echo "Get Failed"
	echo "Returned body: $body"
else
	echo "Passed"
fi

# 3rd test case
echo "3: Testing  (patch json/AF_NEF_PATCH_01.json)"
patch json/AF_NEF_PATCH_01.json $ret_id

if [[ $status_code -ne 200 ]]; then
	echo "Patch Failed"
	echo "Returned body: $body"
else
	echo "Passed"
fi

# 4th test case
echo "4: Testing (delete $ret_id)"
delete $ret_id

if [[ $status_code -ne 204 ]]; then
	echo "Deletion Failed"
	echo "Returned body: $body"
else
	echo "Passed"
fi
