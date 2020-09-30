/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/**
 * @file  PutRequestDispatcher.h
 * @brief Header file for PutRequestDispatcher
 ********************************************************************/

#ifndef __OAMAGENT__PATCHREQUESTDISPATCHER__
#define __OAMAGENT__PATCHREQUESTDISPATCHER__

#include "PatchRequestHandler.h"
#include <stdio.h>
#include <json/json.h>
#include <map>

using namespace std;

class PatchRequestDispatcher
{
    map<string, PatchRequestHandler *> requestHandlers;
public:
    /**
     * @brief        Dispatches a PATCH-method request holding JSON-formatted request data to corresponding handler(s).
     * @param[in]    action      A target handler's designation, as part of the URL.
     * @param[in]    request     JSON-formatted request data.
     * @param[out]   response    A JSON-formatted key-value pair indicating the response.
     * @param[out]   headers     Header of response.
     * @param[in]    cookies     Cookies in request.
     * @throw        Exception   Thrown on failure.
     * @return       void
     */
    void dispatchRequest(const string &action, Json::Value &request, Json::Value &response,
                         map<string, string> &headers, map<string, string> &cookies);
    /**
     * @brief        Links a handler to a URL.
     * @param[in]    action      A designation to assign to a handler, so that requests may be dispatched 
                                 to the handler using that designation as a part of requests' URL.
     * @param[in]    handler     An handler to be linked to the action.
     * @return       void
     */
    void registerHandler(const string &action, PatchRequestHandler &handler);
};

#endif /* defined(__OAMAGENT__PATCHREQUESTDISPATCHER__) */
