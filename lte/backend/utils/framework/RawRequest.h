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
 * @file RawRequest.h
 * @brief Header file for RawRequest
 ********************************************************************/

#ifndef __OAMAGENT__RAWREQUEST__
#define __OAMAGENT__RAWREQUEST__

#include "PostRequestDispatcher.h"
#include "GetRequestDispatcher.h"
#ifdef PUT_SUPPORT
#include "PutRequestDispatcher.h"
#endif
#include "DelRequestDispatcher.h"
#include "PatchRequestDispatcher.h"

#include <fcgio.h>
#include <stdio.h>

class RawRequest
{
    void postRequest(const string &action, FCGX_Request &request, map<string, string> &cookies, PostRequestDispatcher &dispatcher);
    void getRequest(const string &action, FCGX_Request &request, map<string, string> &cookies, GetRequestDispatcher &dispatcher);
    #ifdef PUT_SUPPORT
    void putRequest(const string &action, FCGX_Request &request, map<string, string> &cookies, PutRequestDispatcher &dispatcher);
    #endif
    void delRequest(const string &action, FCGX_Request &request, map<string, string> &cookies, DelRequestDispatcher &dispatcher);
    void patchRequest(const string &action, FCGX_Request &request, map<string, string> &cookies, PatchRequestDispatcher &dispatcher);	
    void printHeaders(map<string, string> &headers);
    const string baseURI;
public:

    /**
     * @brief        Constructs a RawRequest object using a common base URI for the system.
     * @param[in]    baseURI A common base URI for directing requests within the system.
     */
    RawRequest(string &&baseURI);

    /**
     * @brief        Dispatches a FastCGI request.
     * @param[in]    request      A FastCGI request using either POST or GET method.
     * @throw        Exception    Thrown on failure.
     * @return       void
     */
    void dispatch(FCGX_Request &request);

    PostRequestDispatcher postDispatcher;
    GetRequestDispatcher getDispatcher;
    #ifdef PUT_SUPPORT
    PutRequestDispatcher putDispatcher;
    #endif
    DelRequestDispatcher delDispatcher;
    PatchRequestDispatcher patchDispatcher;	
};

#endif /* defined(__OAMAGENT__RAWREQUEST__) */
