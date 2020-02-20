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
#include <string>
#include <json/json.h>
#include <iostream>
#include <fstream>

#include "Log.h"
#include "LocalConfig.h"

using namespace std;

/* ------------------------------------------------------------------------- */
/* Constants */
/* ------------------------------------------------------------------------- */

/* Local configuration variables  */
string localcfg_pgw_ipaddress  = "192.168.120.219";
string localcfg_pgw_port       = "10000";
string localcfg_sgw_ipaddress  = "192.168.120.220";
string localcfg_sgw_port       = "10000";

string localcfg_nginx_hostip   = "127.0.0.1"; 
string localcfg_nginx_port     = "8080"; 
string localcfg_nginx_fcgipass = "127.0.0.1:9999";
string localcfg_http2_enabled = "true";
string localcfg_https_enabled = "true";
string localcfg_ssl_cainfo = "/etc/certs/root-ca-cert.pem";


int oamReadCfgJsonFile(void)
{

    Json::Reader reader;
    Json::Value root;

    //read local.cfg file
    OAMAGENT_LOG(INFO, "Start reading localJsonCfgFile!\n");
        
    ifstream readfile("localcfg.json", ios::binary);
    if(!readfile.is_open())
    {
        OAMAGENT_LOG(ERR, "Open localcfg.json file failed!\n");
        return -1;
    }

    if(reader.parse(readfile,root))
    {
        //get cpf interface information
        localcfg_pgw_ipaddress = root["pgw"]["ipaddress"].asString();
        localcfg_pgw_port      = root["pgw"]["port"].asString();
        OAMAGENT_LOG(INFO, " [PGW] ipaddress  = %s\n", localcfg_pgw_ipaddress.c_str());
        OAMAGENT_LOG(INFO, " [PGW] port       = %s\n", localcfg_pgw_port.c_str());
		
        //get cpf interface information
        localcfg_sgw_ipaddress = root["sgw"]["ipaddress"].asString();
        localcfg_sgw_port      = root["sgw"]["port"].asString();
        OAMAGENT_LOG(INFO, " [SGW] ipaddress  = %s\n", localcfg_sgw_ipaddress.c_str());
        OAMAGENT_LOG(INFO, " [SGW] port       = %s\n", localcfg_sgw_port.c_str());

        // get nginx information
        localcfg_nginx_hostip   = root["nginx"]["hostIp"].asString();
        localcfg_nginx_port     = root["nginx"]["hostPort"].asString();
        localcfg_nginx_fcgipass = root["nginx"]["fcgiPass"].asString();
        OAMAGENT_LOG(INFO, " [NGINX] hostip   = %s\n", localcfg_nginx_hostip.c_str());
        OAMAGENT_LOG(INFO, " [NGINX] port     = %s\n", localcfg_nginx_port.c_str());
        OAMAGENT_LOG(INFO, " [NGINX] fcgipass = %s!\n", localcfg_nginx_fcgipass.c_str());

        // get configuration on whether HTTP2 is to be used
        localcfg_http2_enabled   = root["http2_enabled"].asString();
        OAMAGENT_LOG(INFO, " [HTTP2_ENABLED] = %s\n", localcfg_http2_enabled.c_str());

        // get configuration on whether HTTPS is to be used
        localcfg_https_enabled   = root["https_enabled"].asString();
        OAMAGENT_LOG(INFO, " [HTTPS_ENABLED] = %s\n", localcfg_https_enabled.c_str());

        // get SSL CA Certificate Data. This is used in case of http2_enabled is true
        localcfg_ssl_cainfo   = root["ssl_cainfo"].asString();
        OAMAGENT_LOG(INFO, " [SSL_CAINFO] = %s\n", localcfg_ssl_cainfo.c_str());
    }
    else
    {
        OAMAGENT_LOG(ERR, "Json Parse localcfg.json file failed!\n");  
        return -1;
    }

    readfile.close();
    return 0;

}

