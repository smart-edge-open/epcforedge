/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/**
 * @file  PostRequestHandler.h
 * @brief Header file for Raw request handler
 ********************************************************************/

#ifndef __OAMAGENT__PATCHREQUESTHANDLER__
#define __OAMAGENT__PATCHREQUESTHANDLER__

#include <stdio.h>
#include <json/json.h>
#include <map>

using namespace std;

class PatchRequestHandler
{
public:
    /**
     * @brief        Virtual function to be implemented by REST API Request handlers for processing POST requests.
     * @param[in]    request     JSON-formatted request data.
     * @param[out]   response    A JSON-formatted key-value pair indicating the response.
     * @param[out]   headers     Header of response.
     * @param[in]    cookies     Cookies in request.
     * @throw        Exception   Thrown on failure.
     * @return       void
     */
    virtual void execute(Json::Value &request, Json::Value &response, map<string, string> &headers, map<string, string> &cookies) = 0;
};

#endif /* defined(__OAMAGENT__PATCHREQUESTHANDLER__) */
