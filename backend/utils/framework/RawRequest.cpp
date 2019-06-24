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
 * @file RawRequest.cpp
 * @brief Implementation of RawRequest
 ********************************************************************/

#include "RawRequest.h"
#include "Exception.h"
#include <sstream>
#include <boost/algorithm/string.hpp>
#include "Log.h"

using namespace boost::algorithm;
using namespace std;


RawRequest::RawRequest(string &&baseURI)
  : baseURI(baseURI)
{
//This is a constructor
}

void RawRequest::printHeaders(map<string, string> &headers)
{
    map<string, string>::iterator iter = headers.find("Status");
    if(headers.end() == iter) {
        cout << "Status: 200 OK\r\nContent-type: application/json\r\n";
    } else {
        cout <<"Content-type: application/json\r\n";
    }
    #if 0
    if (headers.size() > 0) {
        pair<string, string> header("",""); 
        for (header : headers) {
           cout << header.first << ": " << header.second << "\r\n";
        }
    }
    #endif
    if (headers.size() > 0) {
        //pair<string, string> header("",""); 
        for (const auto &entry : headers) {
           cout << entry.first << ": " << entry.second << "\r\n";
        }
    }
    cout << "\r\n";
}

void RawRequest::postRequest(const string &action,
                             FCGX_Request &request,
                             map<string, string> &cookies,
                             PostRequestDispatcher &dispatcher)
{
    Json::Value requestJson;
    cin >> requestJson;

    Json::Value responseJson;
    map<string, string> headers;
    dispatcher.dispatchRequest(action, requestJson, responseJson, headers, cookies);
    printHeaders(headers);
    cout << responseJson;
}

void RawRequest::getRequest(const string &action,
                            FCGX_Request &request,
                            map<string, string> &cookies,
                            GetRequestDispatcher &dispatcher)
{
    istringstream queryString(FCGX_GetParam("QUERY_STRING", request.envp));
    OAMAGENT_LOG(INFO, "Processing GetRequest Params.\n");
    string key, val;
    map<string, string> params;
    while (getline(queryString, key, '=') && getline(queryString, val, '&')) {
        params[key] = val;
        OAMAGENT_LOG(INFO, "    Params[%s] = %s.\n", key.c_str(), val.c_str());		
    }

    Json::Value responseJson;
    map<string, string> headers;
    dispatcher.dispatchRequest(action, params, responseJson, headers, cookies);
    printHeaders(headers);
    cout << responseJson;
}
#ifdef PUT_SUPPORT
void RawRequest::putRequest(const string &action,
                             FCGX_Request &request,
                             map<string, string> &cookies,
                             PutRequestDispatcher &dispatcher)
{
    Json::Value requestJson;
    cin >> requestJson;

    Json::Value responseJson;
    map<string, string> headers;
    dispatcher.dispatchRequest(action, requestJson, responseJson, headers, cookies);
    printHeaders(headers);
    cout << responseJson;
}
#endif
void RawRequest::delRequest(const string &action,
                             FCGX_Request &request,
                             map<string, string> &cookies,
                             DelRequestDispatcher &dispatcher)
{
    istringstream queryString(FCGX_GetParam("QUERY_STRING", request.envp));

    string key, val;
    map<string, string> params;
    while (getline(queryString, key, '=') && getline(queryString, val, '&')) {
        params[key] = val;
    }

    Json::Value responseJson;
    map<string, string> headers;
    dispatcher.dispatchRequest(action, params, responseJson, headers, cookies);
    printHeaders(headers);
    cout << responseJson;
}

void RawRequest::patchRequest(const string &action,
                             FCGX_Request &request,
                             map<string, string> &cookies,
                             PatchRequestDispatcher &dispatcher)
{
#if 0
    istringstream queryString(FCGX_GetParam("QUERY_STRING", request.envp));

    string key, val;
    map<string, string> params;
    while (getline(queryString, key, '=') && getline(queryString, val, '&')) {
        params[key] = val;
    }

    Json::Value responseJson;
    map<string, string> headers;
    dispatcher.dispatchRequest(action, params, responseJson, headers, cookies);
    printHeaders(headers);
    cout << responseJson;
#endif
    Json::Value requestJson;
    cin >> requestJson;

    Json::Value responseJson;
    map<string, string> headers;
    dispatcher.dispatchRequest(action, requestJson, responseJson, headers, cookies);
    printHeaders(headers);
    cout << responseJson;	
}


void RawRequest::dispatch(FCGX_Request &request)
{
    string contentType(FCGX_GetParam("CONTENT_TYPE", request.envp));
    string requestMethod(FCGX_GetParam("REQUEST_METHOD", request.envp));
    string documentURI = FCGX_GetParam("DOCUMENT_URI", request.envp);
    istringstream cookies(FCGX_GetParam("HTTP_COOKIE", request.envp) ? FCGX_GetParam("HTTP_COOKIE", request.envp) : "");
    string source_ip(FCGX_GetParam("REMOTE_ADDR", request.envp));
	
    OAMAGENT_LOG(INFO, "Processing FCGX_Request with REMOTE_ADDR %s DocumentURI %s requestMethod %s.\n", 
		 source_ip.c_str(), documentURI.c_str(), requestMethod.c_str());
	
    /* disable the usage of memanager from outside */
    /* We need to disable this check for containers. The nginx server is on different ip and this statement disable memanager ability to work from NES container. */
    /*if (source_ip != me_manager_client_ip &&
        std::find(me_manager_uris.begin(), me_manager_uris.end(), documentURI) != std::end(me_manager_uris)) {
            OAMAGENT_LOG(ERR, "Not allowed to use MeManager from %s\n", source_ip.c_str());
            throw Exception(Exception::DISPATCH_NOTYPE, "Unauthorized MeManager usage");
    }*/

    string key, val;
    map<string, string> cookieParams;
    while (getline(cookies, key, '=') && getline(cookies, val, ';')) {
        trim(val);
        trim(key);
        cookieParams[key] = val;
    }
    OAMAGENT_LOG(INFO, "Get cookie with key (%s) and val (%s) .\n", key.c_str(), val.c_str());


    if ((0 == requestMethod.compare("POST")) && (0 == contentType.compare("application/json"))) {
        OAMAGENT_LOG(INFO, "Calling PostRequest with baseURI (%s) length (%lu).\n", baseURI.c_str(), baseURI.length());
        RawRequest::postRequest(documentURI.substr(baseURI.length()), request, cookieParams, postDispatcher);
		
    } else if (0 == requestMethod.compare("GET")) {
        OAMAGENT_LOG(INFO, "Calling GetRequest with baseURI (%s) length (%lu).\n", baseURI.c_str(), baseURI.length());    
        RawRequest::getRequest(documentURI.substr(baseURI.length()), request, cookieParams, getDispatcher);
		
    } 
#ifdef PUT_SUPPORT
	else if ((0 == requestMethod.compare("PUT")) && (0 == contentType.compare("application/json"))) {
	OAMAGENT_LOG(INFO, "Calling PutRequest with baseURI (%s) length (%lu).\n", baseURI.c_str(), baseURI.length());    
        RawRequest::putRequest(documentURI.substr(baseURI.length()), request, cookieParams, putDispatcher);
		
    } 
#endif
	else if (0 == requestMethod.compare("DELETE")) {
        OAMAGENT_LOG(INFO, "Calling DeleRequest with baseURI (%s) length (%lu).\n", baseURI.c_str(), baseURI.length());   
        RawRequest::delRequest(documentURI.substr(baseURI.length()), request, cookieParams, delDispatcher);
    
    } else if (0 == requestMethod.compare("PATCH")) {
        OAMAGENT_LOG(INFO, "Calling DeleRequest with baseURI (%s) length (%lu).\n", baseURI.c_str(), baseURI.length());   
        RawRequest::patchRequest(documentURI.substr(baseURI.length()), request, cookieParams, patchDispatcher);	
    } else {
        stringstream ss;
        ss << "Method: " << requestMethod << "Content type: " << contentType << " not supported";
        OAMAGENT_LOG(ERR, "%s .\n", ss.str().c_str());
        throw Exception(Exception::DISPATCH_NOTYPE, ss.str());
    }
}
