/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/**
 * @file  DelRequestDispatcher.h
 * @brief Header file for DelRequestDispatcher
 ********************************************************************/

#ifndef __MECFCGI__DELREQUESTDISPATCHER__
#define __MECFCGI__DELREQUESTDISPATCHER__

#include "DelRequestHandler.h"
#include <stdio.h>
#include <json/json.h>
#include <map>

using namespace std;

class DelRequestDispatcher
{
    map<string, DelRequestHandler *> requestHandlers;
public:
    /**
     * @brief        Dispatches a DEL-method request to corresponding handler(s).
     * @param[in]    action      A target handler's designation, as part of the URL.
     * @param[in]    params      Parameters in the DEL request.
     * @param[out]   response    A JSON-formatted key-value pair indicating the response.
     * @param[out]   headers     Header of response.
     * @param[in]    cookies     Cookies in request.
     * @throw        Exception   Thrown on failure.
     * @return       void
     */
    void dispatchRequest(const string &action, map<string, string> &params, Json::Value &response,
                         map<string, string> &headers, map<string, string> &cookies);
    /**
     * @brief        Links a handler to a URL.
     * @param[in]    action      A designation to assign to a handler, so that requests may be dispatched to the handler
                                 using that designation as a part of requests' URL.
     * @param[in]    handler     An handler to be linked to the action.
     * @return       void
     */
    void registerHandler(const string &action, DelRequestHandler &handler);
};

#endif /* defined(__MECFCGI__DELREQUESTDISPATCHER__) */
