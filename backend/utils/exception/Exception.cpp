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
 * @file Exception.cpp
 * @brief Implementation of Exception
 ********************************************************************/

#include "Exception.h"

void Exception::handlerException(Exception e, string &res, string &statusCode)
{
    switch (e.code)
    {
    // ADDED_USERPLANE
    case Exception::ADDED_USERPLANE:
        statusCode = HTTP_SC_ADDED_USERPLANE;
        res= "ADDED_USERPLANE";
        break;	    
    case Exception::INVALID_TYPE:
        statusCode = HTTP_SC_BAD_REQUEST;
        res= "ParameterInvalid";
        break;		
    case Exception::INVALID_PARAMETER:
        statusCode = HTTP_SC_BAD_REQUEST;
        res= "ParameterInvalid";
        break;
    case Exception::INVALID_UERPLANE_FUNCTION:
        res= "INVALID_UERPLANE_PROPERTISE";
        statusCode = HTTP_SC_INVALID_UERPLANE_PROPERTISE;
        break;
    case Exception::INVALID_DATA_SCHEMA:
    case Exception::USERPLANE_NOT_FOUND:
        statusCode = HTTP_SC_USERPLANE_NOT_FOUND;
        res = "USERPLANE_NOT_FOUND";
        break;
    case Exception::INTERNAL_SOFTWARE_ERROR:
        statusCode = HTTP_SC_INTERNAL_SOFTWARE_ERROR;
        res = "INTERNAL_SOFTWARE_ERROR";
        break;    
    case Exception::CONNECT_EPC_ERROR:
        statusCode = HTTP_SC_EPC_CONNECT_ERROR;
        res= "EPC CP Connect failure";
        break;
    case Exception::DISPATCH_NOTARGET:
        statusCode = HTTP_SC_NOT_FOUND;
        res= "404 not found";
        break;
    case Exception::DISPATCH_NOTYPE:
        statusCode = HTTP_SC_BAD_REQUEST;
        res= "BadRequest";
        break;
    default:
        statusCode = HTTP_SC_INTERNAL_SERVER_ERROR;
        res= "UnknownError";
        break;
    }
}
