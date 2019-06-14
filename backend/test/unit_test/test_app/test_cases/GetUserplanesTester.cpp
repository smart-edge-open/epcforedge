/********************************************************************
 * SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2010-2014 Intel Corporation
 ********************************************************************/
/*******************************************************************************
* Integration Tests for AppLiveIndicator, which is a handler for POST requests
* with a payload in JSON.
*******************************************************************************/
#include <json/json.h>

#include "TestUtility.h"
#include "Exception.h"

#include "GetUserplanesTester.h"


const string GET_URL  = {"/userplanes"};
const string GET_URL_ID5  = {"/userplanes/5"};
const string GET_URL_ID10 = {"/userplanes/10"};
const string GET_URL_INVALIDID  = {"/userplanes/99"};
const string GET_URL_INVALID_URI  = {"/userplanes/5/xxx"};



int GetUserplanesTester::execute(string& additionalMessage) {

    try {
        Json::Value resp_body;
        string status_code, discarded_cookie;

        std::cout << "\r\n[RUN       ] Tests for GetUserplanesTester" << std::endl;

        //TesterBase::printState();

        // GET User plane all
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
        sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_GETALL_SUCCESS",
                           HTTP_SC_OK, "OK",
                           status_code, resp_body["result"]);


        // GET User plane ID 5 (DEC format)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
        sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_ID5, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_GETONE_SUCCESS",
                           HTTP_SC_OK, "OK",
                           status_code, resp_body["result"]);

        
	// GET User plane ID INVALID (DEC format)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
        sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_INVALIDID, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_ID_NOTFOUND",
                           HTTP_SC_USERPLANE_NOT_FOUND, "OK",
                           status_code, resp_body["result"]);


        // GET User plane ID 5 (ID_NOTMATCH)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_ID5, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_ID_NOTMATCH",
                           HTTP_SC_INTERNAL_SOFTWARE_ERROR, "OK",
                           status_code, resp_body["result"]);

        // GET User plane ID 5 (TAC_NOTMATCH)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_ID5, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_TAC_NOTMATCH",
                           HTTP_SC_INTERNAL_SOFTWARE_ERROR, "OK",
                           status_code, resp_body["result"]);


        // GET User plane ID 5 (TAC_NOTFOUND)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_ID5, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_TAC_NOTFOUND",
                           HTTP_SC_INTERNAL_SOFTWARE_ERROR, "OK",
                           status_code, resp_body["result"]);


        // GET User plane ID 5 (APNI WRONG)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_ID5, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_APNI_WRONG",
                           HTTP_SC_INTERNAL_SOFTWARE_ERROR, "OK",
                           status_code, resp_body["result"]);


        // GET User plane ID 5 (GetUserPlanesTest_NOSGWS5UWRONG)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_ID5, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_NOSGWS5U",
                           HTTP_SC_INTERNAL_SOFTWARE_ERROR, "OK",
                           status_code, resp_body["result"]);


        // GET User plane ID 5 (GetUserPlanesTest_NOSGWTAC)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_ID5, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_NOSGWTAC",
                           HTTP_SC_INTERNAL_SOFTWARE_ERROR, "OK",
                           status_code, resp_body["result"]);

        // GET User plane ID 5 (GetUserPlanesTest_NOSGWS1U)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_ID5, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_NOSGWS1U",
                           HTTP_SC_INTERNAL_SOFTWARE_ERROR, "OK",
                           status_code, resp_body["result"]);



		// Last MOCK test: negtive test, need mock in curl
        // GET User plane ID all (curl return)
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_EPC_CONN_ERROR",
                           HTTP_SC_EPC_CONNECT_ERROR, "OK",
                           status_code, resp_body["result"]);


        /////////////////////FcgiBackend::run return////////////////////////////////////////
        /////////////////////Need not MOCK function ////////////////////////////////////////
        // GET User plane ID 5 (URI invalid) 
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
        sendGETRequest(status_code, resp_body, TesterBase::cookies["null"], GET_URL_INVALID_URI, "" );
        disconnect();
		reportTestResult("GetUserPlanesTest_INVALIDURI",
                           HTTP_SC_NOT_FOUND, "OK",
                           status_code, resp_body["result"]);

		
        return 0;
    } catch(Exception &e) {
        additionalMessage = e.err;
        disconnect();
        return -1;
    }
}
