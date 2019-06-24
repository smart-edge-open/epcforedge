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
 * @file  PostRequestHandler.h
 * @brief Header file for Raw request handler
 ********************************************************************/

#ifndef __OAMAGENT__JSONREQUESTHANDLER__
#define __OAMAGENT__JSONREQUESTHANDLER__

#include <stdio.h>
#include <json/json.h>
#include <map>

using namespace std;

class PostRequestHandler
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

#endif /* defined(__OAMAGENT__JSONREQUESTHANDLER__) */
