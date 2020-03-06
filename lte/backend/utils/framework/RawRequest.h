/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

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
