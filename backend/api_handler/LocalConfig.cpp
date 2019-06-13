
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

    }
    else
    {
        OAMAGENT_LOG(ERR, "Json Parse localcfg.json file failed!\n");  
        return -1;
    }

    readfile.close();
    return 0;

}

