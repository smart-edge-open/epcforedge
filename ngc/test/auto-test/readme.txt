#SPDX-License-Identifier: Apache-2.0
#Copyright © 2020 Intel Corporation

Installation:
-------------
This test suite requires curl with http2 support. Curl package in yum repo 
doesn't support http2. We need to build curl from source. The script
install_curl.sh automatically download and build curl 7.68.0 which support
http2. nghttp2 is an additional package required by curl and it is installed
automatically by the script. This script requires sudo permissions in order to
work.

Configuration:
--------------
To configure the test framework put the appropriate values in the config file.
This config file is a simple shell script which exports configured variable.
The variables which needed to be configured are:

 1. http_port: Port number on which AF/NEF http server is listening.
 2. https_port: Port number on which AF/NEF https server is listening.
 3. curl_dir: Curl installation directory, it should contain whole path like
	 /home/ahameed/workspace/curl-7.68.0/src/curl
 4. nef_host: The hostname of IP address of AF/NEF server.
 5. https: Set this to true if https should be used for sending requests.
 6. subs_url: This is the subscription path which is the url after port number. 
	[http://nef_host:http_port/subs_url] 


Writing test case:
------------------
The script test_api.sh provide a generic function "send_req" which is used to 
build and send requests based on the arguments passed. To write a test case 
first import the function from test_nef script by adding a line on top of test 
script
"source test_api.sh".

The send_req function takes 4 arguments as
	send_req <method> <sub_id> <body> <expected_response_code>

        method : post/get/put/patch/delete ...
        sub_id : The subscription id retured by the NEF server.
	body : The filepath to json which need to be sended in request body.
	expected_response_code : The http response code expected from the 
		recieved response (200,201,400 ....).

On success send_req function returns 0, on failure it returns non zero value.
The send_req function sets some global variables based on response it got from
the server. These variables can be used by test case writer to validate the
response. The exposed variables are:
	sub_id : This is valid in case of successfull post request. It is the
		subscriber id it got from server.
	status_code : The http response code recieved from server.
	body : The response body recieved from server.

All the 4 arguments to the function are necessary, however for a particular 
request some of the argument may contains dummy data. e.g. for post request
sub_id is not required, however to call send_req any dummy sub_id can be put.

[See test_cases/test_case_2.sh as a sample test case]


Running:
--------

Before running test script make sure:
 1. Curl with http2 support is installed.
 2. Appropriate values are setted in config file.
 3. AF/NEF server is reachable from testing machine [Some time proxy might cause
	issue]

To run a test script just execute it
Ex: [user@abc]./test_cases/test_case_2.sh


Known issues:
-------------
While executing ./install_curl.sh some failure may occur due to some package 
not installed. To build curl some packages are required, by default they are 
installed on basic centos. If any issue occurs try installing these packages
and run install_curl.sh again.

Extra packages that might be required:
 g++ make binutils autoconf automake autotools-dev libtool pkg-config 
 zlib1g-dev libcunit1-dev libssl-dev libxml2-dev libev-dev libevent-dev
 libjansson-dev libjemalloc-dev cython python3-dev python-setuptools

Some time proxy might cause some issue, check proxy setting before running
install_curl.sh.

