#!/bin/bash

#SPDX-License-Identifier: Apache-2.0
#Copyright © 2020 Intel Corporation

# Test Case Scenario 
# 1) POST 001 with app1,app2
# 2) GET 001,app1
# 3) PUT 001,app1 with updated pfd1,pfd2
# 4) PUT 001,app1 with invalid pfd1
# 5) PATCH 001,app1 with changed pfd2
# 6) PATCH 001,app1 with invalid pfd2
# 7) DELETE 001,app1
# 8) GET 001,app1 (deleted app)
# 9) DELETE 001,app10
# 10) DELETE trans 001

source ../../auto-test/test_api.sh
source ../../auto-test/test_pfd.sh
 
appId=0
echo -e "\n\n\\tAPPLICATION PATCH/GET/PUT/DELETE TESTS"
echo -e "\t-----------------------------------------------"

# calling send_req function for post request with json file path
# <json/AF_NEF_PFD_POST_001.json> and expected response http code is 201.
echo -e "\t PFD POST TRANS: \c" 
if send_req post 0 json/AF_NEF_PFD_POST_001.json 201; then
	get_trans_id
        get_app_id_in_trans
        passed 
        echo -e "\t\t TRANSACTION ID: " $trans_id
        echo -e "\t\t APPLICATIONS: " $app_arr_trans

else
	failed $status_code
fi

# Assigning the app_id to the first element of app_arr
appId=`echo $app_arr_trans | awk '{print $1}'`

# Calling send_req function for get request with the trans_id and app_id 
# previously returned from the post request. Here json file path is dummy as 
# get request doesn't require json file path. Expected response is 200
echo -e "\t PFD GET APP $appId in $trans_id\c"
if send_app_req get $trans_id $appId dummy_json 200; then
        passed
else
        failed $status_code
fi


# Calling send_req function for put request with the trans_id  and  
# app_id previously Expected response is 200
echo -e "\t PFD APP PUT $appId in $trans_id: \c"
if send_app_req put $trans_id $appId json/AF_NEF_PFD_APP_PUT_001.json 200; then
        passed
else
        failed $status_code
fi

# Calling send_req function for put request with the trans_id  and  
# app_id previously from POST. Invalid PFD. Expected response is 200
echo -e "\t PFD APP INVALID PUT $appId in $trans_id: \c"
if send_app_req put $trans_id $appId json/AF_NEF_PFD_APP_PUT_002.json 400; then
        passed
else
        failed $status_code
fi

# Calling send_req function for patch request with the trans_id  and  
# app_id previously Expected response is 200
echo -e "\t PFD APP PATCH $appId in $trans_id : \c" 
if send_app_req patch $trans_id $appId json/AF_NEF_PFD_APP_PATCH_001.json 200; \
then
        passed
else
        failed $status_code
fi

# Calling send_req function for patch request with the trans_id  and  
# app_id previously from POST. Invalid PFD. Expected response is 200
echo -e "\t PFD APP INVALID PATCH $appId in $trans_id: \c"
if send_app_req patch $trans_id $appId json/AF_NEF_PFD_APP_PATCH_002.json 400; \
then
        passed
else
        failed $status_code
fi

# Calling send_req function for delete request with the trans_id  and 
#  app_id previously 
# Here json file path is dummy as delete request 
# doesn't require json file path. Expected response is 204
echo -e "\t PFD APP DELETE $appId in $trans_id: \c" 
if send_app_req delete $trans_id $appId dummy_json 204; then
        passed
else
        failed $status_code
fi


# Calling send_req function for get request with the trans_id  and invalid 
# app_id previously Here json file path is dummy as get request 
# doesn't require json file path. Expected response is 200
echo -e "\t PFD APP GET INVALID APP $appId in $trans_id: \c" 
if send_app_req get $trans_id $app_id dummy_json 404; then
        passed 
else
        failed $status_code
fi

# Calling send_req function for delete request with the trans_id  and invalid
# app_id. Here json file path is dummy as delete request 
# Expected response is 404
echo -e "\t PFD APP INVALID DELETE app10 in $trans_id: \c" 
if send_app_req delete $trans_id app10 dummy_json 404; then
        passed
else
        failed $status_code
fi

# Calling send_req function for get request with the trans_id  and invalid 
# app_id previously Here json file path is dummy as get request 
# doesn't require json file path. Expected response is 200
echo -e "\t PFD APP DELETE TRANS $trans_id: \c" 
if send_req delete $trans_id dummy_json 204; then
        passed
else
        failed $status_code
fi


