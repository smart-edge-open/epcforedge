/*******************************************************************************
* Integration Tests for AppLiveIndicator, which is a handler for POST requests
* with a payload in JSON.
*******************************************************************************/
#include <json/json.h>

#include "TestUtility.h"
#include "Exception.h"

#include "PatchUserplanesTester.h"


const string PATCH_URL = {"/userplanes/5"};
const string PATCH_URL_INVALID_URI  = {"/userplanes/5/xxx"};


int PatchUserplanesTester::execute(string& additionalMessage) {

    try {
        Json::Value resp_body;
        string status_code, discarded_cookie;

        std::cout << "\r\n[RUN       ] Tests for PatchUserplanesTest" << std::endl;
        //TesterBase::printState();

        // SUCCESS
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
		sendPATCHRequest(status_code, resp_body, TesterBase::cookies["null"], PATCH_URL, "" , "PatchUserplanes");
        disconnect();
		reportTestResult("PatchUserPlanesTest_SUCCESS", 
			                        HTTP_SC_OK, "OK", 
			                        status_code, resp_body["result"]);
        #if 1
        // HTTP_SC_ADDED_USERPLANE
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
		sendPATCHRequest(status_code, resp_body, TesterBase::cookies["null"], PATCH_URL, "" , "PatchUserplanes");
        disconnect();
		reportTestResult("PatchUserPlanesTest_NOTFOUND", 
			                        HTTP_SC_USERPLANE_NOT_FOUND, "OK", 
			                        status_code, resp_body["result"]);

        #endif
     
        // negtive testing, curl return failured
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
		sendPATCHRequest(status_code, resp_body, TesterBase::cookies["null"], PATCH_URL, "" , "PatchUserplanes");
        disconnect();
		reportTestResult("PatchUserPlanesTest_EPC_CONN_ERROR", 
			                        HTTP_SC_EPC_CONNECT_ERROR, "OK", 
			                        status_code, resp_body["result"]);


        /////////////////////FcgiBackend::run return////////////////////////////////////////
        /////////////////////Need not MOCK function ////////////////////////////////////////
        // PATCH User plane ID 5 (URI invalid) 
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPATCHRequest(status_code, resp_body, TesterBase::cookies["null"], PATCH_URL_INVALID_URI,"", "PatchUserplanes" );
        disconnect();
		reportTestResult("PatchUserPlanesTest_INVALIDURI",
                           HTTP_SC_NOT_FOUND, "OK",
                           status_code, resp_body["result"]);

        return 0;
    } catch(Exception &e) {
        additionalMessage = e.err;
        disconnect();
        return -1;
    }
}
