/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/**
 * @file  GetRequestHandler.h
 * @brief Header file for Raw request handler
 ********************************************************************/

#ifndef __OAMAGENT__GETREQUESTHANDLER__
#define __OAMAGENT__GETREQUESTHANDLER__

#include <stdio.h>
#include <json/json.h>
#include <map>

using namespace std;

class GetRequestHandler
{
public:
    /**
     * @brief        Virtual function to be implemented by REST API Request handlers for processing GET requests.
     * @param[in]    params      Parameters in a GET request.
     * @param[out]   response    A JSON-formatted key-value pair indicating the response.
     * @param[out]   headers     Header of response.
     * @param[in]    cookies     Cookies in request.
     * @throw        Exception   Thrown on failure.
     * @return       void
     */
    virtual void execute(map<string, string> params, Json::Value &response, map<string, string> &headers,
                         map<string, string> &cookies) = 0;
};

#endif /* defined(__OAMAGENT__GETREQUESTHANDLER__) */
