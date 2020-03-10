/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/**
 * @file  PatchRequestDispatcher.cpp
 * @brief PATCH method and JSON formatted request dispatcher.
 ********************************************************************/

#include "PatchRequestDispatcher.h"
#include "Exception.h"
#include "Log.h"

#define REQUEST_SPLIT_MARK         "/"
#define REQUEST_SPLIT_MARK_LENGTH  1

/**
 * @brief		 Dispatches a PATCH-method request holding JSON-formatted request data to corresponding handler(s).
 * @param[in]	 action 	 A target handler's designation, as part of the URL.
 * @param[in]	 request	 JSON-formatted request data.
 * @param[out]	 response	 A JSON-formatted key-value pair indicating the response.
 * @param[out]	 headers	 Header of response.
 * @param[in]	 cookies	 Cookies in request.
 * @throw		 Exception	 Thrown on failure.
 * @return		 void
 */
void PatchRequestDispatcher::dispatchRequest(const string &action,
                                            Json::Value &request,
                                            Json::Value &response,
                                            map<string, string> &headers,
                                            map<string, string> &cookies)
{

	OAMAGENT_LOG(INFO, "PatchRequestDispatcher with action %s.\n", action.c_str()); 	


    size_t pos;
    if (requestHandlers.find(action) != requestHandlers.end()) {
        static_cast<PatchRequestHandler *>(requestHandlers[action])->execute(request, response, headers, cookies);
        return;
    } else if ((pos = action.find_last_of(REQUEST_SPLIT_MARK)) != string::npos) {
        string newAction = action.substr(0, pos) + "/UUID";
        if (requestHandlers.find(newAction) != requestHandlers.end()) {
            request["UUID"] = action.substr(pos + REQUEST_SPLIT_MARK_LENGTH, action.length() - pos + REQUEST_SPLIT_MARK_LENGTH);
            static_cast<PatchRequestHandler *>(requestHandlers[newAction])->execute(request, response, headers, cookies);
            return;
        }
    }
    OAMAGENT_LOG(ERR, "Dispatch failed, action: %s.\n", action.c_str());
    throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");

	
}

/**
 * @brief        Links a handler to a URL.
 * @param[in]	 action 	 A designation to assign to a handler, so that requests may be dispatched to the handler using that designation as a part of requests' URL.
 * @param[in]	 handler	 An handler to be linked to the action.
 * @return       void
 */
void PatchRequestDispatcher::registerHandler(const string &action, PatchRequestHandler &handler)
{
    requestHandlers[action] = &handler;
}
