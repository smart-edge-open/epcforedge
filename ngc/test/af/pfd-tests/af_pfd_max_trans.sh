#!/bin/bash
#SPDX-License-Identifier: Apache-2.0
#Copyright © 2020 Intel Corporation


# Test Case Scenario 
# 1) PFD POST 11 transactions
# 2) PFD GET ALL 
# 3) PFD DELETE <all transactions>

# Sets the config  and includes all the lib functions
source ../../auto-test/test_api.sh
source ../../auto-test/test_pfd.sh

num=1
echo -e "\n\n\tPFD MAX TRANSACTIONS TESTS"
echo -e "\t-----------------------------------------------"

echo -e "\t CREATE PFD MAX TRANS: "
while [[ $num -ne 12 ]]; do
        str="app${num}"
        sed -i "s/app.*[^\"]\"/${str}\"/g" json/AF_NEF_PFD_POST_004.json
        # Sending post request untill it starts to fail
        if send_req post 0 json/AF_NEF_PFD_POST_004.json 201; then
                get_trans_id
                echo -e "\t\c"
                passed "${transId:?}"
        else

		if [[ $(status_code) == 400 ]]; then
                	echo -e "\t\tMAX TRANS REACHED, POST FAILED"
                	echo -e "\t\c"
			passed
		else
			failed "${status_code:?}"
             	fi
        fi
        (( num++ ))
done


# calling send_req function for get_all request with dummy json 
# expected response 200 (no PFDs)
echo -e "\t PFD GET ALL: \c"
if send_req get_all 0 dummy_json 200; then
        get_all_trans
        passed
        echo -e "\t\t TRANSACTION ID: ${trans_arr:?}"

else
        failed "${status_code:?}"
fi

num=1

echo -e "\t DELETE PFD MAX TRANS: "
while [[ $num -ne 11 ]]; do
        # Sending delete requests for all trans
        transId=$(echo "$trans_arr" | awk -v i="$num" '{print $i}')
        if send_req delete "$transId" dummy_json 204; then
                echo -e "\t\c"
                passed "${transId:?}"
        else
                echo -e "\t\c"
                failed "${status_code:?}"

        fi
        (( num++ ))
done

