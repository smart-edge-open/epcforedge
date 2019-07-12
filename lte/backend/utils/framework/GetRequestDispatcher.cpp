/*******************************************************************************
* Copyright 2019 Intel Corporation. All rights reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*******************************************************************************/
/**
 * @file  GetRequestDispatcher.cpp
 * @brief GET method requests dispatcher
 ********************************************************************/

#include "GetRequestDispatcher.h"
#include "Exception.h"
#include "Log.h"

#define REQUEST_SPLIT_MARK         "/"
#define REQUEST_SPLIT_MARK_LENGTH  1
void GetRequestDispatcher::dispatchRequest(const string &action,
                                            map<string, string> &params,
                                            Json::Value &response,
                                            map<string, string> &headers,
                                            map<string, string> &cookies)
{
    OAMAGENT_LOG(INFO, "GetRequestDispatcher with action: %s.\n", action.c_str()); 	

    if (!action.length()) {
        OAMAGENT_LOG(ERR, "Dispatch failed.\n");
        throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
    }
    //string action_bk = action;
    size_t pos;
    if (requestHandlers.find(action) != requestHandlers.end()) {
    	OAMAGENT_LOG(INFO, "GetRequestDispatcher Find execute handler for the action (%s).\n", action.c_str()); 			
        static_cast<GetRequestHandler *>(requestHandlers[action])->execute(params, response, headers, cookies);
        return;
    } else if ((pos = action.find_last_of(REQUEST_SPLIT_MARK)) != string::npos) {
    	OAMAGENT_LOG(INFO, "GetRequestDispatcher Find SplitMask for the action %s.\n", action.c_str());
        string newAction = action.substr(0, pos) + "/UUID";
        if (requestHandlers.find(newAction) != requestHandlers.end()) {
    	    //OAMAGENT_LOG(INFO, "GetRequestDispatcher substr(%d,%d) with pos =%d, actlen=%d\n",
            params["UUID"] = action.substr(pos + REQUEST_SPLIT_MARK_LENGTH, action.length() - pos + REQUEST_SPLIT_MARK_LENGTH);
    	    OAMAGENT_LOG(INFO, "GetRequestDispatcher Find id (%s) for the newaction (%s)\n",params["UUID"].c_str(), newAction.c_str());
	    if (0 == strlen(params["UUID"].c_str())) {
                throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
            }
            static_cast<GetRequestHandler *>(requestHandlers[newAction])->execute(params, response, headers, cookies);
            return;
			
        }
    }
    OAMAGENT_LOG(ERR, "Dispatch failed, action: %s.\n", action.c_str());
    throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
}

void GetRequestDispatcher::registerHandler(const string &action, GetRequestHandler &handler)
{
    requestHandlers[action] = &handler;
}
