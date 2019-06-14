/********************************************************************
 * SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2010-2014 Intel Corporation
 ********************************************************************/
/**
 * @file main.cpp
 * @brief Implementation of Main for OAM agent
 ********************************************************************/

#include <sys/types.h>
#include <unistd.h>
#include <signal.h>
#include "Log.h"
#include "FcgiBackend.h"
#include "RawRequest.h"
#include "UpfController.h"
#include "LocalConfig.h"


using namespace std;

int init(int argc, char *argv[])
{

    // oamagent: log function
    oamagentLogInit();

    // Load config file
    if (0 != oamReadCfgJsonFile()) {
       OAMAGENT_LOG(ERR, "JsonCfgFile Read Failed, exit!\n");
       return -1;
    }
    
    // Init FastCGI backend
    //string hostIp      = "127.0.0.1";
    //string hostPort    = "8080";
    //string fastcgiPass = "127.0.0.1:9999";
    FcgiBackend::setNginxInfo(localcfg_nginx_hostip, uint16_t(atoi((char *)localcfg_nginx_port.c_str())), localcfg_nginx_fcgipass);
	
    return 0;
    
}

/**
* @brief			The function will register control APIs
*					
* @param[in]		raw 	  The raw request.
* @return			void.
*/
void registerUpfControlAPI(RawRequest &raw)
{

    UserplaneAdd       addUserplane; //POST
    UserplanePatchByID patchUserplaneByID; //PATCH
    UserplanesListGet  getUseplanesList; //GET
    UserplaneGetByID   getUserplaneByID; //GET
    UserplaneDelByID   delUserplaneByID; //DELETE

	//GET userplanes
    raw.getDispatcher.registerHandler("userplanes",  getUseplanesList);
	//POST userplanes
    raw.postDispatcher.registerHandler("userplanes",  addUserplane);
	//GET userplanes/UUID
    raw.getDispatcher.registerHandler("userplanes/UUID",  getUserplaneByID);
    //DELETE userplanes/UUID
    raw.delDispatcher.registerHandler("userplanes/UUID",  delUserplaneByID);
    //PATCH userplanes/UUID
    raw.patchDispatcher.registerHandler("userplanes/UUID",  patchUserplaneByID);


}

/**
* @brief			The function will handle signalss
*					
* @param[in]		signal 	  The signal.
* @return			void.
*/
void handle_signals(int signal) {
    if (SIGTERM == signal) {
        OAMAGENT_LOG(INFO, "SIGTERM SIGINT, exiting now\n");
        exit(0);
    }
}


/**
* @brief			The OAM agent main entry
*					
* @param[in]		argc 	  The arg num
* @param[in]		argv 	  The arg values
* @return			error code.
*/

int main(int argc, char *argv[])
{
    // oamagent: register signals handler
    OAMAGENT_LOG(INFO, "OAMAgent Starting ........\n");
    struct sigaction act;
    act.sa_handler = &handle_signals;
    act.sa_flags = 0;
    act.sa_flags |= SA_RESTART;
    sigfillset(&act.sa_mask);
    if (sigaction(SIGTERM, &act, NULL) == -1) {
        cerr<<"Cannot handle SIGTERM"<<endl;
        return -1;
    }

    if (0 != init(argc, argv)){
		OAMAGENT_LOG(ERR, "Init failed, existing... \n");    
        return -1;
    }

    // oamagent: register fcgi handler 
    OAMAGENT_LOG(INFO, "OAMAgent Backend MgmtAPI registering ....\n");    
    RawRequest raw ("/");
    registerUpfControlAPI(raw);

	// oamagent: run fcgi backend 
    OAMAGENT_LOG(INFO, "OAMAgent Backend running ....\n");    	
    FcgiBackend::run(raw);

    return 0;

}
