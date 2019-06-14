/********************************************************************
 * SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2010-2014 Intel Corporation
 ********************************************************************/
/**
 * @file    UpfController.h
 * @brief   Classes of handlers responsible for handling UPF control message
 *          from MEC controller
 ********************************************************************/

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
