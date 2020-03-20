#SPDX-License-Identifier: Apache-2.0
#Copyright © 2020 Intel Corporation

#!/bin/bash


# Before running the scripts, the pre-requisite is to run the AF and NEF 
# on baremetal with same config as in the config file

# In the test scripts, the transaction id is extracted from the self link
# in the response body. Similarly , the applicaion id is extracted from self 
# link of application. The individual test scenario files list the test cases

count_pass=0
count_fail=0

# Sets the config  and includes all the lib functions
source ../../auto-test/test_api.sh
source ../../auto-test/test_pfd.sh


# Tests for transaction GET/POST/PUT/DELETE
. ./af_pfd_trans_tests.sh

# Tests for application GET/PUT/PATCH/DELETE
. ./af_pfd_app_tests.sh

# Test for Max transactions
. ./af_pfd_max_trans.sh

# Display the summary of tests
echo -e "\n\n\t SUMMARY OF TESTS"
echo -e "\t-----------------------------------------"
display_summary
