#!/bin/bash
#SPDX-License-Identifier: Apache-2.0
#Copyright © 2020 Intel Corporation


# A sample test case to show json unmarshalling using jq tool. body variable 
# after making post request contains json data. jq tool is used to unmarshall
# the json data. To get value belong to particular key (e.g. self)
#	echo $body | jq -r '.self'


# Sets the config  and includes the lib function
source test_nef.sh

# calling send_req function for post request with json file path
# <json/AF_NEF_POST_02.json> and expected response http code is 201.
if send_req post 0 json/AF_NEF_POST_02.json 201; then
	# on success post request $body contains returned json data
	echo "The returned response body is:"
	echo "------------------------------"
	echo "${body:?}"
	echo "------------------------------"

	echo "The returned json data after unmarshalling is"
	echo "---------------------------------------------"
	echo "${body:?}" | jq # Print the returned body in json format.

	echo -n "The value belong to key ipv4Addr is: "
	echo "${body:?}" | jq -r '.ipv4Addr' # Get the value belong to key (ipv4Addr) in json data.
else
	echo "Failed"
fi

