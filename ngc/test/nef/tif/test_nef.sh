#!/bin/bash

#SPDX-License-Identifier: Apache-2.0
#Copyright � 2020 Intel Corporation

# Script to build & send http/https request to server and validate the response.
# Before calling any function of the script make sure appropriate values are set
# in the config file.


# Default configuration
http_port=8061
https_port=8060
#curl_path=/home/ahameed/workspace/curl-7.68.0/src/curl
curl_path=curl
sub_url=3gpp-traffic-influence/v1/AF_01/subscriptions
sub_url1=3gpp-traffic-influence/v1
sub_url2=subscriptions
nef_host=localhost
https=false
cert_path=/etc/certs/root-ca-cert.pem

# Import the configuration variable from config file.
configure()
{
	source config
}

# Print current configuration.
show_config()
{
	echo "HTTP Port: $http_port"
	echo "HTTPS Port: $https_port"
	echo "Curl Path: $curl_path"
	echo "Subscriber URI: $sub_url"
	echo "NEF Host: $nef_host"
	echo "Using HTTPS: $https" 
}

# Sends post request to the configured server. It required one argument, json 
# file path (which needs to be send in request body).
# e.g. post json/AF_NEF_POST_02.json
# On successfull execution it returns 0 and set sub_id to the returned 
# subscription id.
post()
{
	if [[ ! -f $1 ]]; then
		echo "Invalid Json filepath"
		return 1
	fi
	
	if [[ $https == "true" ]]; then
		out=`$curl_path -w '\nResponse Status=%{http_code}\n' --cacert $cert_path --http2 -X \
			POST -H "Content-Type: application/json" --data @$1 \
			https://$nef_host:$https_port/$sub_url 2>/dev/null`
	else
		 out=`$curl_path -w '\nResponse Status=%{http_code}\n' -X POST -H \
		 "Content-Type: application/json" --data @$1 \
		 http://$nef_host:$http_port/$sub_url 2>/dev/null`

		
	fi
	echo "post req body sent" 
	jq . $1
	status_code=`echo $out | grep "Response Status"  | \
	awk 'BEGIN { FS="=" } // {print $2}'`
	if [[ $status_code == 201 ]]; then
		body=`echo $out | sed 's/ Response Status.*//g'`
		sub_id=`echo $body | jq -r '.self' | awk 'BEGIN {FS="/"} \
		// {print $(NF)}'`
		echo "post res body received" 
		echo $body | jq .
	fi
	return 0
}


# Sends get request to the configured server. It required one argument, sub_id
# e.g. get <sub_id>
# On successfull execution it returns 0.
get()
{
	sub_id=$1
	if [[ $sub_id =~ ^[0-9].+$ ]]; then
		if [[ $https == "true" ]]; then
			out=`$curl_path  -w '\nResponse Status=%{http_code}\n' --cacert $cert_path \
			--http2 -X GET https://$nef_host:$https_port/$sub_url/$sub_id \
			2>/dev/null`
		else
			 out=`$curl_path  -w '\nResponse Status=%{http_code}\n' -X \
			 GET http://$nef_host:$http_port/$sub_url/$sub_id 2>/dev/null`
			
		fi
			
		status_code=`echo $out | grep "Response Status"  | awk \
		'BEGIN { FS="=" } // {print $2}'`
		body=`echo $out | sed 's/ Response Status.*//g'`
		echo "res body received" 
		echo $body | jq .
	else
		echo "Invalid sub_id"
		return 2
	fi
	return 0
}
# Sends get request to the configured server.
# e.g. get_all
# On successfull execution it returns 0.
get_all()
{
	
	
	if [[ $https == "true" ]]; then
		out=`$curl_path  -w '\nResponse Status=%{http_code}\n' --cacert $cert_path \
		--http2 -X GET https://$nef_host:$https_port/$sub_url 2>/dev/null`
	else
		 out=`$curl_path  -w '\nResponse Status=%{http_code}\n' -X GET \
		 http://$nef_host:$http_port/$sub_url 2>/dev/null`
		
	fi
	status_code=`echo $out | grep "Response Status"  | awk 'BEGIN \
	{ FS="=" } // {print $2}'`
	body=`echo $out | sed 's/ Response Status.*//g'`
	return 0
}

# Sends put request to the configured server. It required two arguments, json 
# file path and sub_id.
# e.g. put <json_file_path> <sub_id>
# On successfull execution it returns 0.
put()
{
	sub_id=$2
	if [[ ! -f $1 ]]; then
		echo "Invalid Json filepath"
		return 1
	fi

	if [[ $sub_id =~ ^[0-9].+$ ]]; then
		if [[ $https == "true" ]]; then
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' --cacert $cert_path \
			--http2 -X PUT -H "Content-Type: application/json" --data @$1 \
			https://$nef_host:$https_port/$sub_url/$sub_id 2>/dev/null`
		else
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' -X PUT -H \
			"Content-Type: application/json" --data @$1 \
			http://$nef_host:$http_port/$sub_url/$sub_id 2>/dev/null`
		fi
		status_code=`echo $out | grep "Response Status"  | awk 'BEGIN \
		{ FS="=" } // {print $2}'`
		body=`echo $out | sed 's/ Response Status.*//g'`
		echo "res body sent" 
		echo $body | jq .
	else
		echo "Invalid sub_id"
		return 2
	fi
	return 0

}

# Sends patch request to the configured server. It required two arguments, json 
# file path and sub_id.
# e.g. put <json_file_path> <sub_id>
# On successfull execution it returns 0.
patch()
{
	if [[ ! -f $1 ]]; then
		echo "Invalid Json filepath"
		return 1
	fi

	sub_id=$2
		if [[ $sub_id =~ ^[0-9].+$ ]]; then
		if [[ $https == "true" ]]; then
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' --cacert $cert_path --http2 \
			-X PATCH -H "Content-Type: application/json" --data @$1 \
			https://$nef_host:$https_port/$sub_url/$sub_id 2>/dev/null`
		else
			 out=`$curl_path -w '\nResponse Status=%{http_code}\n' -X PATCH -H \
			 "Content-Type: application/json" --data @$1 \
			 http://$nef_host:$http_port/$sub_url/$sub_id 2>/dev/null`

		fi
		echo "patch req body sent" 
	    jq . $1
		status_code=`echo $out | grep "Response Status"  | awk 'BEGIN \
		{ FS="=" } // {print $2}'`
		body=`echo $out | sed 's/ Response Status.*//g'`
		echo "patch res body received" 
		echo $body | jq .
	else
		echo "Invalid sub_id"
		return 2
	fi
	return 0
}

# Sends delete request to the configured server. It required one arguments 
# sub_id.
# e.g. delete <sub_id>
# On successfull execution it returns 0.
delete()
{
	sub_id=$1
	if [[ $sub_id =~ ^[0-9].+$ ]]; then
		if [[ $https == "true" ]]; then
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' --cacert $cert_path --http2 \
			-X DELETE https://$nef_host:$https_port/$sub_url/$sub_id 2>/dev/null`
		else
			out=`$curl_path -w '\nResponse Status=%{http_code}\n' -X DELETE \
			http://$nef_host:$http_port/$sub_url/$sub_id 2>/dev/null`
	
		fi
		status_code=`echo $out | grep "Response Status"  | awk 'BEGIN \
		{ FS="=" } // {print $2}'`
		body=`echo $out | sed 's/ Response Status.*//g'`
		
	else
		echo "Invalid sub_id"
		return 2
	fi
	return 0
}

# Build and send request to the configured servers based on the arguments. This
# function also validate the expeceted response and return 0 if returned http
# response code match expected response code.
# Usage:
# 	send_req <method> <sub_id> <data> <expected_response>
#
# All the 4 arguments are compulsory, in case an argument is not required put 
# some dummy data. e.g for post request sub_id is not required, so the function
# call can have any dummy sub_id 
#    send_req post 001(dummy sub_id) <data> <expected_response>
send_req()
{
	method=$1
	sub_id=$2
	data=$3
	expected_response=$4
	ret_val=false
	case "$method" in
		"delete")
			delete $sub_id
			;;
		"get")
			get $sub_id
			;;
		"patch")
			patch $data $sub_id
			;;
		"post")
			post $data
			;;
		"put")
			put $data
			;;
		"get_all")
			get_all
			;;
		*)
			echo "Invalid Method"
			return 2
			;;
	esac
echo $status_code
	if [[ $status_code -ne $expected_response ]]; then
		ret_val=false
		return 1
	else
		ret_val=true
		return 0
	fi

}

# Configure the framework variables.
#configure