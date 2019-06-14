/********************************************************************
 * SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2010-2014 Intel Corporation
 ********************************************************************/
/**
 * @file  DelRequestHandler.h
 * @brief Header file for DEL request handler
 ********************************************************************/

#ifndef __OAMAGENT__DELREQUESTHANDLER__
#define __OAMAGENT__DELREQUESTHANDLER__

#include <stdio.h>
#include <json/json.h>
#include <map>

using namespace std;

class DelRequestHandler
{
public:
    /**
     * @brief        Virtual function to be implemented by REST API Request handlers for processing DEL requests.
     * @param[in]    params      Parameters in a DEL request.
     * @param[out]   response    A JSON-formatted key-value pair indicating the response.
     * @param[out]   headers     Header of response.
     * @param[in]    cookies     Cookies in request.
     * @throw        Exception   Thrown on failure.
     * @return       void
     */
    virtual void execute(map<string, string> params, Json::Value &response, map<string, string> &headers, map<string, string> &cookies) = 0;
};

#endif /* defined(__OAMAGNET__DELREQUESTHANDLER__) */
