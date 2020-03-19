#!/bin/bash

#SPDX-License-Identifier: Apache-2.0
#Copyright © 2019 Intel Corporation


# A sample test case to test post, get, patch and delete request to server. 


# Sets the config  and includes the lib function
source test_nef.sh


#exposed env variables

# 1. $sub_id  //Valid after successfull POST req
# 2. $status_code // Valid for all methods when server respods
# 3. $body // body in the response received 
# 

# calling send_req function for post request with json file path
# <json/AF_NEF_POST_02.json> and expected response http code is 201.
if send_req post 0 json/AF_NEF_POST_02.json 201; then
	echo "Passed" $sub_id # sub_id variable is set by the send_req func. 
else
	echo "Failed"
fi


# Calling send_req function for get request with the sub_id previously returned
# from the post request. Here json file path is dummy as get request doesn't
# require json file path.
if send_req get $sub_id dummy_json 200; then
	echo "Passed"
else
	echo "Failed"
fi

# Calling send_req function for patch request with the sub_id previously
# returned from the post request. 
if send_req patch $sub_id json/AF_NEF_PATCH_01.json 200; then
	echo "Passed"
else
	echo "Failed"
fi

# Calling send_req function for delete request with the sub_id previously
# returned from the post request. 
if send_req delete $sub_id dummy 204; then
	echo "Passed"
else
	echo "Failed"
fi
