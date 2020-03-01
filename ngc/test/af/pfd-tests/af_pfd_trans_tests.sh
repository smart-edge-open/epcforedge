#!/bin/bash

# Test Case Scenarios in this script 
# 1) PFD GET ALL (no PFD) 
# 2) PFD POST with two apps - app1, app2
# 3) PFD POST with same apps - app1, app2
# 4) PFD POST with one valid app and one duplicate
# 5) PFD POST with incorrect json
# 6) PFD GET <trans_id> ( Returned in the Self link from POST )
# 7) PFD GET <invalid_trans>
# 8) PFD GET ALL
# 9) PFD PUT <trans_id>
# 10) PFD PUT with invalid pfds
# 11) PFD DELETE
# 12) PFD DELETE <invalid_trans>
# 13) PFD GET <trans_id>

echo -e "\n\n\tTRANSACTION POST/GET/PUT/DELETE TESTS"
echo -e "\t-----------------------------------------------"

# calling send_req function for get_all request with dummy json 
# expected response 200 (no PFDs)
echo -e "\t PFD GET ALL: \c"
if send_req get_all 0 dummy_json 200; then
        passed
        
else
        failed $status_code
fi


# calling send_req function for post request with json file path
# <json/AF_NEF_PFD_POST_001.json> and expected response http code is 201.
echo -e "\t PFD POST TRANS: \c"
if send_req post 0 json/AF_NEF_PFD_POST_001.json 201; then
	get_trans_id
        get_app_id_in_trans
        trans_id_1=$trans_id  #trans id extracted from self
	passed
        echo -e "\t\t TRANSACTION ID: " $trans_id_1
        echo -e "\t\t APPLICATIONS: " $app_arr_trans # extracted from self
else
	failed $status_code
fi


# calling send_req function for post request with same json file path
# <json/AF_NEF_PFD_POST_001.json> and expected response http code is 
# 500(Duplicate appID)
echo -e "\t PFD DUPILCATE POST: \c" 
if send_req post 0 json/AF_NEF_PFD_POST_001.json 500; then
        passed
else
        failed $status_code
fi

# calling send_req function for post request with json file path
# <json/AF_NEF_PFD_POST_002.json>. This has one duplicate app and one valid app
# expected response http code is 200 ( PFD report for duplicate app)
echo -e "\t PFD POST (ONE APP INVALID): \c"
if send_req post 0 json/AF_NEF_PFD_POST_003.json 201; then
        get_trans_id
        get_app_id_in_trans
        trans_id_2=$trans_id
        passed
        echo -e "\t\t TRANSACTION ID: " $trans_id_2
        echo -e "\t\t APPLICATIONS: " $app_arr_trans # extracted from self
else
        failed $status_code
fi

# calling send_req function for post request with json file path
# <json/AF_NEF_PFD_POST_007.json>. This has json error
# expected response http code is 400 ( BAD request)
echo -e "\t PFD POST INVALID JSON: \c"
if send_req post 0 json/AF_NEF_PFD_POST_007.json 400; then
        passed
else
        failed $status_code
fi

# Calling send_req function for get request with the trans_id previously 
# returned from the post request. Here json file path is dummy as get request 
# doesn't require json file path. Expected response is 200
echo -e "\t PFD GET TRANS $trans_id_1: \c"
if send_req get $trans_id_1 dummy_json 200; then
	passed $trans_id1
else
	failed $status_code
fi

# Calling send_req function for get request with invalid transaction id
# Here json file path is dummy as get request doesn't require json file path.
# expected response is 404
echo -e "\t PFD GET INVALID TRANS 20000: \c"
if send_req get 20000 dummy_json 404; then
        passed
else
        failed $status_code
fi


# calling send_req function for get_all request with dummy json 
# expected response 200 (PFDs)
echo -e "\t PFD GET ALL: \c" 
if send_req get_all 0 dummy_json 200; then
        passed
        get_all_trans
        get_all_app
        echo -e "\t\t TRANSACTION IDs: " $trans_arr
        echo -e "\t\t APPLICATION IDs: " $app_arr
       
else
        failed $status_code
fi

# calling send_req function for PUT request with json (updated PFDs)
# expected response 200 (PFDs)
echo -e "\t PFD PUT TRANS $trans_id_1: \c"
if send_req put $trans_id_1 json/AF_NEF_PFD_PUT_001.json 200; then
        passed
else
        failed $status_code
fi

# calling send_req function for PUT request with json (updated PFDs)
# expected response 200 (PFDs)
echo -e "\t PFD INVALID PUT $trans_id_1: \c"
if send_req put $trans_id_1 json/AF_NEF_PFD_PUT_002.json 400; then
        passed
else
        failed $status_code
fi


# Calling send_req function for delete request with the trans_id previously
# returned from the post request. Expected response is 204
echo -e "\t PFD DELETE TRANS $trans_id_1: \c"
if send_req delete $trans_id_1 dummy 204; then
	passed
else
	failed $status_code
fi

# Calling send_req function for delete request with the trans_id previously
# returned from the post request. Expected response is 204
echo -e "\t PFD DELETE TRANS $trans_id_2: \c"
if send_req delete $trans_id_2 dummy 204; then
	passed
else
	failed $status_code
fi

# Calling send_req function for delete request with invalid transaction id
# Expected response code is 404
echo -e "\t PFD INVALID DELETE TRANS 20000: \c"
if send_req delete 20000 dummy 404; then
        passed
else
        failed $status_code
fi

# Calling send_req function for get request with the trans_id previously 
# returned from the post request. Here json file path is dummy as get request 
# doesn't require json file path. Expected response is 404 as trabs_id is 
# deleted
echo -e "\t PFD GET PFD TRANS(DELETED) $trans_id_1: \c"
if send_req get $trans_id_1 dummy_json 404; then
	passed
else
	failed $status_code
fi


