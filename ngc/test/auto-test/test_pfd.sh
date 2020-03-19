#!/bin/bash

#SPDX-License-Identifier: Apache-2.0
#Copyright © 2019 Intel Corporation

		
# Script to build & send http/https request to server for PFD applications and 
# validate the resp nse.
# Before calling any function of the script make sure appropriate values are set
# in the config file.
trans_id=0
# GET ALL - array of PFD trans, get trans ids from self
get_all_trans()
{
	trans_arr=`echo $body | jq '.[] | .self' | awk -F "/" \
        '{print $NF}' | sed 's/"//g'`
	
}

# GET-  get trans id from self in a single PFD trans
get_trans_id()
{
        trans_id=`echo $body | jq .self | awk -F "/" \
        '{print $NF}' | sed 's/"//g'`
        
}

# GET ALL APP IDs -array of PFD trans, get all app ids from self in PFD Data
get_all_app()
{
 	app_arr=`echo $body | jq '.[] | .pfdDatas[].self' | awk -F "/" \
         '{print $NF}' | sed 's/"//g'`
        
	
}

# GET ALL APP IDs - get all app ids from self in a single PFD Data
get_app_id_in_trans()
{
        app_arr_trans=`echo $body | jq .pfdDatas[].self | awk -F "/" \
         '{print $NF}' | sed 's/"//g'`
        
}

# Printing the pass in colored format
passed()
{
       echo -e  "\e[1;32m PASSED \e[0m" $1
       count_pass=$((count_pass+1))
}

# Printing the fail in colored format
failed()
{
       echo -e  "\e[1;31m FAILED \e[0m" $1
       count_fail=$((count_fail+1))
}

# Display of the summary of results
display_summary()
{
        echo -e "\e[1;32m\t TOTAL PASSED = \e[0m" $count_pass
        echo -e "\e[1;31m\t TOTAL FAILED = \e[0m" $count_fail
}


# Sends get request to the configured server for an application.
# It required two argument, trans_id and app_id
# e.g. get <trans_id> <app_id>
# On successfull execution it returns 0.
get_app()
{
	trans_id=$1
	app_id=$2
	if [[ $trans_id =~ ^[0-9].+$ ]]; then
		if [[ $https == "true" ]]; then
			out=`$curl_path  -w '\nResponse Status=%{http_code}\n' --cacert $cert_path --http2 \
			-X GET https://$nef_host:$https_port/$sub_url/$trans_id/applications/$app_id \
			2>/dev/null`
		else
			out=`$curl_path  -w '\nResponse Status=%{http_code}\n' -X GET \
	http://$nef_host:$http_port/$sub_url/$trans_id/applications/$app_id \
	2>/dev/null`
		fi
		status_code=`echo $out | grep "Response Status"  | \
		awk 'BEGIN { FS="=" } // {print $2}'`
		body=`echo $out | sed 's/ Response Status.*//g'`
	else
		echo "Invalid trans_id"
		return 2
	fi
	return 0
}

# Sends put request to the configured server for an application. 
#It required three arguments, json 
# file path, trans_id and app_id
# e.g. put <json_file_path> <trans_id> <app_id>
# On successfull execution it returns 0.
put_app()
{
	trans_id=$2
	app_id=$3
	if [[ ! -f $1 ]]; then
		echo "Invalid Json filepath"
		return 1
	fi

	if [[ $trans_id =~ ^[0-9].+$ ]]; then
		if [[ $https == "true" ]]; then
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' --cacert $cert_path --http2 \
			-X PUT -H "Content-Type: application/json" --data @$1 \
		https://$nef_host:$https_port/$sub_url/$trans_id/applications/$app_id \
			2>/dev/null`
		else
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' -X PUT -H \
			"Content-Type: application/json" --data @$1 \
		http://$nef_host:$http_port/$sub_url/$trans_id/applications/$app_id \
			2>/dev/null`
		fi
		status_code=`echo $out | grep "Response Status"  | awk \
		'BEGIN { FS="=" } // {print $2}'`
		body=`echo $out | sed 's/ Response Status.*//g'`
	else
		echo "Invalid trans_id"
		return 2
	fi
	return 0

}

# Sends patch request to the configured server for an application.
# It required three arguments, json 
# file path, trans_id and app_id
# e.g. put <json_file_path> <trans_id> <app_id>
# On successfull execution it returns 0.
patch_app()
{
	if [[ ! -f $1 ]]; then
		echo "Invalid Json filepath"
		return 1
	fi

	trans_id=$2
	app_id=$3
	if [[ $trans_id =~ ^[0-9].+$ ]]; then
		if [[ $https == "true" ]]; then
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' --cacert $cert_path --http2 \
			-X PATCH -H "Content-Type: application/json" --data @$1 \
		https://$nef_host:$https_port/$sub_url/$trans_id/applications/$app_id \
			2>/dev/null`
		else
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' -X PATCH -H \
			"Content-Type: application/json" --data @$1 \
		http://$nef_host:$http_port/$sub_url/$trans_id/applications/$app_id \
		2>/dev/null`
		fi
		status_code=`echo $out | grep "Response Status"  | \
		awk 'BEGIN { FS="=" } // {print $2}'`
		body=`echo $out | sed 's/ Response Status.*//g'`
	else
		echo "Invalid trans_id"
		return 2
	fi
	return 0
}

# Sends delete request to the configured server for an application. It 
#required two arguments trans_id and appId
# e.g. delete <trans_id> <app_id>
# On successfull execution it returns 0.
delete_app()
{
	trans_id=$1
	app_id=$2
	if [[ $trans_id =~ ^[0-9].+$ ]]; then
		if [[ $https == "true" ]]; then
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' --cacert $cert_path --http2 \
			-X DELETE \
		https://$nef_host:$https_port/$sub_url/$trans_id/applications/$app_id \
			2>/dev/null`
		else
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' -X DELETE \
		http://$nef_host:$http_port/$sub_url/$trans_id/applications/$app_id \
			2>/dev/null`
		fi
		status_code=`echo $out | grep "Response Status"  | \
		awk 'BEGIN { FS="=" } // {print $2}'`
		body=`echo $out | sed 's/ Response Status.*//g'`
	else
		echo "Invalid trans_id"
		return 2
	fi
	return 0
}

# Build and send request to the configured servers based on the arguments. This
# function also validate the expeceted response and return 0 if returned http
# response code match expected response code.
# Usage:
# 	send_req <method> <trans_id> <app_id> <data> <expected_response>
#
# All the 5 arguments are compulsory, in case an argument is not required put 
# some dummy data.  
#    send_app_req get <trans_id> <app_id> <data> <expected_response>
send_app_req()
{
	method=$1
	trans_id=$2
	app_id=$3
	data=$4
	expected_response=$5
	ret_val=false
	case "$method" in
		"delete")
			delete_app $trans_id $app_id
			;;
		"get")
			get_app $trans_id $app_id
			;;
		"patch")
			patch_app $data $trans_id $app_id
			;;
		"put")
			put_app $data $trans_id $app_id
			;;
		*)
			echo "Invalid Method"
			return 2
			;;
	esac

	if [[ $status_code -ne $expected_response ]]; then
		ret_val=false
		return 1
	else
		ret_val=true
		return 0
	fi

}

# Configure the framework variables.
configure
