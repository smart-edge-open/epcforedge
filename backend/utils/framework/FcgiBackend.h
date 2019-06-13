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
* @file  FcgiBackend.h
* @brief FastCGI Backend manager for OAM agent.
*/

#ifndef __OAMAGENT__FCGIBACKEND__
#define __OAMAGENT__FCGIBACKEND__

#include <iostream>
#include "RawRequest.h"

using namespace std;

/* ------------------------------------------------------------------------- */
/* Macros */
/* ------------------------------------------------------------------------- */

//#define FCGI_BUF_SIZE 256
#define FCGI_BUF_SIZE 81920


/* ------------------------------------------------------------------------- */
/* Class Defs */
/* ------------------------------------------------------------------------- */

class FcgiBackend
{

private:
    typedef struct nginx_info_s
    {
        string      hostIp;
        uint16_t    hostPort;
        string      fastcgiPass;
    } nginx_info_t;
    static nginx_info_t nginxInfo;

public:
    /**
    * @brief        Configure Nginx related information
    * @param[in]    hostIp       Host IP address.
    * @param[in]    hostPort     Host Port.
    * @param[in]    fastcgiPass  FastCGI Pass.    
    * @return       void
    */
    static void setNginxInfo(string hostIp, uint16_t hostPort, string fastcgiPass);

    /**
    * @brief        Starts processing RESTful API requests.
    * @param[out]    RawRequest     raw request object.
    * @return       void
    */
    static void run(RawRequest &raw);

};

#endif /* __OAMAGENT__FCGIBACKEND__ */
