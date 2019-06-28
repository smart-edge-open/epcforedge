```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation.
```

# api_test

This is the CURL based test scripts for HTTPS API testing


## Overview

- Check whether installed CURL support SSL.
- To run test scripts, need copy pre-generated self-signed certificate file (epc.crt) from "/etc/nginx/ssl/" to this place.
- If OAMAgent and CURL TestScripts will run on the same server, you will need to add hostname "epc.oam" into /etc/hosts 127.0.0.1 entry
- If OAMAgent and CURL TestScripts will run on the different servers, you will need to add hostname "epc.oam" into /etc/hosts on the CURL Test scripts Server. 
  And the IP address should be OAMAgent's IPAddress

