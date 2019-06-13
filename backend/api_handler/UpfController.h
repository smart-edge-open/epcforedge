/*******************************************************************************
Copyright 2019 Intel Corporation. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*******************************************************************************/
/**
 * @file    UpfController.h
 * @brief   Classes of handlers responsible for handling UPF control message
 *          from MEC controller
 *
 */

#ifndef __OAMAGENT__UPFCONTROLLER__
#define __OAMAGENT__UPFCONTROLLER__

#include "PostRequestHandler.h"
#ifdef PUT_SUPPORT
#include "PutRequestHandler.h"
#endif
#include "GetRequestHandler.h"
#include "DelRequestHandler.h"
#include "PatchRequestHandler.h"
#include "CupsMgmtInterface.h"

/* ------------------------------------------------------------------------- */
/* Macros */
/* ------------------------------------------------------------------------- */

#define MAX_USERPLANES_NUM 128


/* ------------------------------------------------------------------------- */
/* Class Defs */
/* ------------------------------------------------------------------------- */

class UserplaneAdd : public PostRequestHandler
{
public:
    CupsMgmtMessage cupsMgmtMsg;
    /**
    * @brief            Post new user plane configuration
    * @param[in]        request     JSON-formatted request data.
    * @param[out]       response    JSON-formatted key-value pair(s) indicating
    *                               response.
    * @param[out]       headers     Response headers.
    * @param[in]        cookies     Cookies header in request.
    * @throw            Exception   Thrown on failure.
    * @return           void
    */

    void execute(Json::Value &request, Json::Value &response,
                 map<string, string> &headers, map<string, string> &cookies);
};

class UserplanePatchByID : public PatchRequestHandler
{
public:
	CupsMgmtMessage cupsMgmtMsg;	
    /**
    * @brief            Update the configuration of specific user plane
    * @param[in]        request     JSON-formatted request data.
    * @param[out]       response    JSON-formatted key-value pair(s) indicating
    *                               response.
    * @param[out]       headers     Response headers.
    * @param[in]        cookies     Cookies header in request.
    * @throw            Exception   Thrown on failure.
    * @return           void
    */

    void execute(Json::Value &request, Json::Value &response,
                 map<string, string> &headers, map<string, string> &cookies);
};

class UserplaneDelByID : public DelRequestHandler
{
public:
	CupsMgmtMessage cupsMgmtMsg;	
    /**
    * @brief            Delete specific user plane.
    * @param[in]        params      JSON-formatted params data.
    * @param[out]       response    JSON-formatted key-value pair(s) indicating
    *                               response.
    * @param[out]       headers     Response headers.
    * @param[in]        cookies     Cookies header in request.
    * @throw            Exception   Thrown on failure.
    * @return           void
    */
    void execute(map<string, string> params, Json::Value &response,
                 map<string, string> &headers, map<string, string> &cookies);
};


class UserplanesListGet : public GetRequestHandler
{
public:
	CupsMgmtMessage cupsMgmtMsg;
	
    /**
    * @brief            Get User Planes List.
    * @param[in]        request     JSON-formatted request data.
    * @param[out]       response    JSON-formatted key-value pair(s) indicating
    *                               response.
    * @param[out]       headers     Response headers.
    * @param[in]        cookies     Cookies header in request.
    * @throw            Exception   Thrown on failure.
    * @return           void
    */
    void execute(map<string, string> params, Json::Value &response,
                 map<string, string> &headers, map<string, string> &cookies);
};

class UserplaneGetByID : public GetRequestHandler
{

public:
	CupsMgmtMessage cupsMgmtMsg;	
	
    /**
    * @brief            Get Specific User Plane.
    * @param[in]        request     JSON-formatted request data.
    * @param[out]       response    JSON-formatted key-value pair(s) indicating
    *                               response.
    * @param[out]       headers     Response headers.
    * @param[in]        cookies     Cookies header in request.
    * @throw            Exception   Thrown on failure.
    * @return           void
    */
    void execute(map<string, string> params, Json::Value &response,
                 map<string, string> &headers, map<string, string> &cookies);
};

#endif //__OAMAGENT__UPFCONTROLLER__