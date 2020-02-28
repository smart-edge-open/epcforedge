/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

#include <json/json.h>

#include "TestUtility.h"
#include "Exception.h"

#include "DelUserplanesTester.h"


const string DEL_URL = {"/userplanes"};
const string DEL_URL_ID5 = {"/userplanes/5"};
const string DEL_URL_ID10 = {"/userplanes/10"};
const string DEL_URL_INVALIDID = {"/userplanes/99"};
const string DEL_URL_INVALIDURI = {"/userplanes/99/xxx"};



int DelUserplanesTester::execute(string& additionalMessage) {

    try {
        Json::Value resp_body;
        string status_code, discarded_cookie;

        std::cout << "\r\n[RUN       ] Tests for DelUserplanesTester" << std::endl;

        //TesterBase::printState();

        // DEL User plane Success
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
        sendDELRequest(status_code, resp_body,discarded_cookie, DEL_URL_ID5, "");
        disconnect();
        reportTestResult("DelUserPlanesTest_SUCCESS", 
             HTTP_SC_OK, "OK", status_code, resp_body["result"]);

        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
        sendDELRequest(status_code, resp_body,discarded_cookie, DEL_URL_INVALIDID, "");
        disconnect();
        reportTestResult("DelUserPlanesTest_IDNOTFOUND", 
               HTTP_SC_USERPLANE_NOT_FOUND, "OK",status_code, resp_body["result"]);

        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
	sendDELRequest(status_code, resp_body,discarded_cookie, DEL_URL_ID5, "");
        disconnect();
	reportTestResult("DelUserPlanesTest_EPC_CONN_ERROR", 
              HTTP_SC_USERPLANE_NOT_FOUND, "OK", status_code, resp_body["result"]);


        /////////////////////FcgiBackend::run return////////////////////////////////////////
        /////////////////////Need not MOCK function ////////////////////////////////////////
        // Del User plane ID 5 (URI invalid) 
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
	sendDELRequest(status_code, resp_body, discarded_cookie, DEL_URL_INVALIDURI, "" );
        disconnect();
	reportTestResult("DelUserPlanesTest_INVALIDURI",
                           HTTP_SC_NOT_FOUND, "OK",
                           status_code, resp_body["result"]);
		

	#if 0
        printf("report 1\n");
        // DEL User plane NOTFOUND
        status_code = "";
        resp_body.clear();	        
	sendDELRequest(status_code, resp_body,discarded_cookie, DEL_URL_INVALIDID, "");
        printf("report 1.1\n");		
        disconnect();
	printf("report 1.2\n");
        reportTestResult("DelUserPlanesTest_ID_NOFOUND", 
			               HTTP_SC_USERPLANE_NOT_FOUND, "USERPLANE_NOT_FOUND",
                           status_code, resp_body["result"]);
        printf("report 2\n");		
		#endif
        return 0;
    } catch(Exception &e) {
        additionalMessage = e.err;
        disconnect();
        return -1;
    }
}
