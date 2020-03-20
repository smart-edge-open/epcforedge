#SPDX-License-Identifier: Apache-2.0
#Copyright © 2020 Intel Corporation

#!/bin/bash

# A sample test case to test the maximum subscription a server can handle. 


# Import the send_req function.
source test_nef.sh


num=1
failed_req=0
while [[ $num -ne 30 ]]; do
	# Sending post request untill it starts to fail
	if send_req post 0 json/AF_NEF_POST_02.json 201; then
		echo "Test $num Passed sub_id=$sub_id"
	else
		(( failed_req++ ))
		echo "Test $num Failed $failed_req time"
	fi

	if [[ $failed_req -eq 10 ]]; then
		break;
	fi
	(( num++ ))
done
