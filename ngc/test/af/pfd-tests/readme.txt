#SPDX-License-Identifier: Apache-2.0
#Copyright © 2020 Intel Corporation

Purpose:
-------------

The tests in this folder are for the functionality testing of AF and NEF
The http requests are sent to AF which inturn sends them to NEF to get the
response.
AF and NEF should be running on baremetal before running these tests.

Installation:
-------------

These tests use the framework which is under ngc/test/auto-test. 
Follow the installation instructions given in the readme.txt under 
ngc/test/auto-test

Configuration:
--------------
To configure the test framework put the appropriate values in the config file.
This config file is a simple shell script which exports configured variables
The variables which needed to be configured are:

 1. http_port: Port number on which AF http server is listening.
 2. https_port: Port number on which AF https server is listening.
 3. curl_dir: Curl installation directory, it should contain whole path like
	 /home/ahameed/workspace/curl-7.68.0/src/curl
 4. nef_host: The hostname of IP address of AF server.
 5. https: Set this to true if https should be used for sending requests.
 6. subs_url: This is the subscription path which is the url after port number. 
	[http://nef_host:http_port/subs_url] 


Running the test cases:
------------------

The folder has the following files:

af_pfd_trans_tests.sh - This has the transaction level tests 
						GET/POST/PUT/DELETE
af_pfd_max_trans.sh - This has the tests for maximum transaction limit in NEF
af_pfd_app_tests.sh - This has all the application Level tests 
						GET/PUT/PATCH/DELETE
af_pfd_all_tests.sh - This is a script to run all the above test scripts


Prereuisites:
 1. Curl with http2 support is installed.
 2. Appropriate values are setted in config file.
 3. AF and NEF executables are running on baremetal

