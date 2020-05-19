#SPDX-License-Identifier: Apache-2.0
#Copyright � 2020 Intel Corporation

#!/bin/bash

source test_nef.sh
#post request with sub_id=0,json file as post body,expected response code(for success case) and AF id  
if send_req post 0 json/AF_NEF_POST_01.json 201; then
	echo "Passed" $sub_id # sub_id variable is set by the send_req func. 
else
	echo "Failed"
fi
#get request with sub_id as received from previous post,dummy json file,expected response code(for success case) and AF id 
if send_req get $sub_id dummy_jsonn 200; then
	echo "Passed" $sub_id # sub_id variable is set by the send_req func. 
else
	echo "Failed"
fi
#patch request with sub_id as received from previous post,json file as patch body,expected response code(for success case) and AF id 
if send_req patch $sub_id json/AF_NEF_PATCH_01.json 200; then
	echo "Passed" $sub_id # sub_id variable is set by the send_req func. 
else
	echo "Failed"
fi
#delete request with sub_id as received from previous post,dummy json file,expected response code(for success case) and AF id 
if send_req delete $sub_id dummy_jsonn 204; then
	echo "Passed" $sub_id # sub_id variable is set by the send_req func. 
else
	echo "Failed"
fi
# if send_req post 0 json/AF_NEF_POST_02.json 201 AF_01; then
# 	echo "Passed" $sub_id # sub_id variable is set by the send_req func. 
# else
# 	echo "Failed"
# fi
# if send_req get_all 0 dummy_jsonn 200 AF_01; then
# 	echo "Passed" $sub_id # sub_id variable is set by the send_req func. 
# else
# 	echo "Failed"
# fi