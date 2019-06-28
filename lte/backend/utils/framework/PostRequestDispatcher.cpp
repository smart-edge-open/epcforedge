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
 * @file  PostRequestDispatcher.cpp
 * @brief POST method and JSON formatted request dispatcher.
 ********************************************************************/

#include "PostRequestDispatcher.h"
#include "Exception.h"
#include "Log.h"

void PostRequestDispatcher::dispatchRequest(const string &action,
                                            Json::Value &request,
                                            Json::Value &response,
                                            map<string, string> &headers,
                                            map<string, string> &cookies)
{

    OAMAGENT_LOG(INFO, "PostRequestDispatcher with action %s.\n", action.c_str()); 
    if (!action.length()) {
        OAMAGENT_LOG(ERR, "Dispatch failed.\n");
        throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
    }

    if (requestHandlers.find(action) != requestHandlers.end()) {
		OAMAGENT_LOG(INFO, "PostRequestDispatcher Find execute handler for the action %s.\n", action.c_str()); 	
        static_cast<PostRequestHandler *>(requestHandlers[action])->execute(request, response, headers, cookies);
        return;
    }
    OAMAGENT_LOG(ERR, "Dispatch failed, action: %s.\n", action.c_str());
    throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
}

void PostRequestDispatcher::registerHandler(const string &action, PostRequestHandler &handler)
{
    requestHandlers[action] = &handler;
}
