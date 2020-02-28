/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/**
 * @file  DelRequestDispatcher.cpp
 * @brief DELETE method requests dispatcher
 *********************************************************************/

#include "DelRequestDispatcher.h"
#include "Exception.h"
#include "Log.h"

#define REQUEST_SPLIT_MARK         "/"
#define REQUEST_SPLIT_MARK_LENGTH  1
void DelRequestDispatcher::dispatchRequest(const string &action,
                                            map<string, string> &params,
                                            Json::Value &response,
                                            map<string, string> &headers,
                                            map<string, string> &cookies)
{

    OAMAGENT_LOG(INFO, "DelRequestDispatcher with action %s.\n", action.c_str()); 	

    if (!action.length()) {
        OAMAGENT_LOG(ERR, "Dispatch failed.\n");
        throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
    }

    size_t pos;
    if (requestHandlers.find(action) != requestHandlers.end()) {
		OAMAGENT_LOG(INFO, "DelRequestDispatcher Find execute handler for the action %s.\n", action.c_str());
        static_cast<DelRequestHandler *>(requestHandlers[action])->execute(params, response, headers, cookies);
        return;
    } else if ((pos = action.find_last_of(REQUEST_SPLIT_MARK)) != string::npos) {
        string newAction = action.substr(0, pos) + "/UUID";
        if (requestHandlers.find(newAction) != requestHandlers.end()) {
            params["UUID"] = action.substr(pos + REQUEST_SPLIT_MARK_LENGTH, action.length() - pos + REQUEST_SPLIT_MARK_LENGTH);
			OAMAGENT_LOG(INFO, "DelRequestDispatcher Find UUID %s for the newaction %s.\n",params["UUID"].c_str(), newAction.c_str());	
            static_cast<DelRequestHandler *>(requestHandlers[newAction])->execute(params, response, headers, cookies);
            return;
        }
    }
    OAMAGENT_LOG(ERR, "Dispatch failed, action: %s.\n", action.c_str());
    throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
}

void DelRequestDispatcher::registerHandler(const string &action, DelRequestHandler &handler)
{
    requestHandlers[action] = &handler;
}
