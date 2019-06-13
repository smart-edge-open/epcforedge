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
* @file  FcgiBackend.cpp
* @brief FastCGI backend manager for OAM Agent.
*/

#include "FcgiBackend.h"
#include "Log.h"
#include <boost/asio.hpp>
#include "Exception.h"

using namespace boost::asio;

FcgiBackend::nginx_info_t FcgiBackend::nginxInfo;
const int FCGX_BACKLOG = 99;

/**
* @brief		Configure Nginx related information
* @param[in]	hostIp		 Host IP address.
* @param[in]	hostPort	 Host Port.
* @param[in]	fastcgiPass  FastCGI Pass.	  
* @return		void
*/
void FcgiBackend::setNginxInfo(string hostIp, uint16_t hostPort, string fastcgiPass)
{
    nginxInfo.hostIp = hostIp;
    nginxInfo.hostPort = hostPort;
    nginxInfo.fastcgiPass = fastcgiPass;
}

/**
* @brief		Starts processing RESTful API requests.
* @param[out]	 RawRequest 	raw request object.
* @return		void
*/
void FcgiBackend::run(RawRequest &raw)
{
    char buf[FCGI_BUF_SIZE];
    FCGX_Request request;
    FCGX_Init();
    dup2(FCGX_OpenSocket(nginxInfo.fastcgiPass.c_str(), FCGX_BACKLOG), 0);
    FCGX_InitRequest(&request, 0, 0);

    while (0 == FCGX_Accept_r(&request)) {
        fcgi_streambuf cin_fcgi_streambuf(request.in, buf, FCGI_BUF_SIZE);
        fcgi_streambuf cout_fcgi_streambuf(request.out, buf, FCGI_BUF_SIZE);
        std::streambuf* in_backup = cin.rdbuf(&cin_fcgi_streambuf);
        std::streambuf* out_backup = cout.rdbuf(&cout_fcgi_streambuf);

        try {
            raw.dispatch(request);
        }
        catch (Exception &e) {
            switch(e.code) {
                case Exception::DISPATCH_NOTARGET:
                    cout << "Status: 404 Not Found\r\n\r\n";
                    break;
                case Exception::DISPATCH_NOTYPE:
                    cout << "Status: 400 Bad Request\r\n\r\n";
                    break;
                case Exception::PARSING_JSON_BODY:
                    cout << "Status: 400 Bad Request\r\n\r\n";
                    break;
                default:
                    cout << "Status: 520 Unknown Error\r\n\r\n";
            }
            OAMAGENT_LOG(ERR, "ERROR(%d):%s \n", e.code, e.err.c_str());
        }
        catch (runtime_error::exception &e) {
            OAMAGENT_LOG(ERR, "%s \n",  e.what());
        }
        cin.rdbuf(in_backup);
        cout.rdbuf(out_backup);
    }
}
