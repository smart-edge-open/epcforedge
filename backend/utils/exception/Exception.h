/********************************************************************
 * SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2010-2014 Intel Corporation
 ********************************************************************/
/**
 * @file Exception.h
 * @brief Header file for Exception
 ********************************************************************/

#ifndef __OAMAGENT__EXCEPTION__
#define __OAMAGENT__EXCEPTION__

#include <stdio.h>
#include <iostream>
#include <map>
#include <json/json.h>

using namespace std;

/* ------------------------------------------------------------------------- */
/* Macros */
/* ------------------------------------------------------------------------- */

#define HTTP_SC_SWITCHING_PROTOCOLS "101 Switching Protocols"
#define HTTP_SC_OK "200 OK. Operation Successful"
#define HTTP_SC_CREATED "201 Created"
#define HTTP_SC_NO_CONTENT "204 No Content"
#define HTTP_SC_BAD_REQUEST "400 Bad Request"
#define HTTP_SC_INVALID_UERPLANE_PROPERTISE "400 Invalid userplane properties provided"
#define HTTP_SC_NOT_FOUND "404 Not Found"
#define HTTP_SC_USERPLANE_NOT_FOUND "404 Userplane not found"
#define HTTP_SC_ADDED_USERPLANE "201 Added Userplane"

#define HTTP_SC_INTERNAL_SERVER_ERROR "500 Internal Server Error"
#define HTTP_SC_EPC_CONNECT_ERROR "500 EPC CP Connect Error"
#define HTTP_SC_INTERNAL_SOFTWARE_ERROR "500 Internal Software Error"

/* ------------------------------------------------------------------------- */
/* Class Defs */
/* ------------------------------------------------------------------------- */

class Exception
{
public:
    enum {
        DISPATCH_NOTARGET,
        DISPATCH_NOTYPE,
        ADDED_USERPLANE,
        INVALID_ACTION,
        USERPLANE_NOT_FOUND,
        INVALID_PARAMETER,
        INTERNAL_SOFTWARE_ERROR,
        INVALID_TYPE,
        INVALID_UERPLANE_FUNCTION,
        CONNECT_EPC_ERROR,
        PARSING_JSON_BODY,
        INVALID_DATA_SCHEMA
    };

    Exception(int code, const string &err) : code(code), err(err) { }
    Exception(int code) : code(code) { }
    int code;
    string err;
    static void handlerException(Exception e, string &res, string &statusCode);
};

#endif /* defined(__OAMAGENT__EXCEPTION__) */
