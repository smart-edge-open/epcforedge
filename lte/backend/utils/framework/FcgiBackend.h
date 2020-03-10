/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/**
 * @file  FcgiBackend.h
 * @brief FastCGI Backend manager for OAM agent.
 ********************************************************************/

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
