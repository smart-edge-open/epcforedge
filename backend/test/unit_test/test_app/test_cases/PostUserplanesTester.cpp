/*******************************************************************************
* Integration Tests for AppLiveIndicator, which is a handler for POST requests
* with a payload in JSON.
*******************************************************************************/
#include <json/json.h>

#include "TestUtility.h"
#include "Exception.h"

#include "PostUserplanesTester.h"


const string POST_URL = {"/userplanes"};
const string POST_URL_INVALIDURI = {"/userplanes/xxx"};


int PostUserplanesTester::execute(string& additionalMessage) {

    try {
        Json::Value resp_body;
        string status_code, discarded_cookie;

        std::cout << "\r\n\[RUN       ] Tests for PostUserplanesTest" << std::endl;
        //TesterBase::printState();

        ///////////////////////CURL MOCK return ////////////////////////////
        // SUCCESS : 2 POST MOCK
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL, "" , "PostUserplanes");
        disconnect();
		reportTestResult("PostUserPlanesTest_ID5_SUCCESS", 
			                        HTTP_SC_OK, "OK", 
			                        status_code, resp_body["result"]);

        // HTTP_SC_ADDED_USERPLANE: actually PGW POST will cause successflag false 1 POST MOCK
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL, "" , "PostUserplanes");
        disconnect();
		reportTestResult("PostUserPlanesTest_PGW ADDED_USERPLANE", 
			                        HTTP_SC_ADDED_USERPLANE, "OK", 
			                        status_code, resp_body["result"]);



        // HTTP_SC_ADDED_USERPLANE: actually SGW POST will cause successflag false 2 POST MOCK
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL, "" , "PostUserplanes");
        disconnect();
		reportTestResult("PostUserPlanesTest_SGW_ADDED_USERPLANE", 
			                        HTTP_SC_ADDED_USERPLANE, "OK", 
			                        status_code, resp_body["result"]);

	    // PGW SGW ID not match
        status_code = ""; 
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL, "" , "PostUserplanes");
        disconnect();
		reportTestResult("PostUserPlanesTest_PGW SGW_IDNOTMATCH", 
			                        HTTP_SC_INTERNAL_SOFTWARE_ERROR, "OK", 
			                        status_code, resp_body["result"]);






		
        /////// last mock function
        /////// Negtive testing. MOCK in the cpfCurlPost
        // CONNECT_EPC_ERROR
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);
		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL, "" , "PostUserplanes");
        disconnect();
		reportTestResult("PostUserPlanesTest_EPC_CONN_ERROR", 
			                        HTTP_SC_EPC_CONNECT_ERROR, "OK", 
			                        status_code, resp_body["result"]);




        /////////////////////BELOW Need not CRUL MOCK function ////////////////////////////////////////
        /////////////////////UserplaneAdd::API Check with wrong json request////////////////////////////////////////
        /////////////////////Need not MOCK function ////////////////////////////////////////
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanesWrongFunction" );
        disconnect();
		reportTestResult("PostUserPlanesTest_INVALIDFUNCTION",
                           HTTP_SC_INVALID_UERPLANE_PROPERTISE, "OK",
                           status_code, resp_body["result"]);
		
        /////////////////////UserplaneAdd::API Check with wrong json request////////////////////////////////////////
        /////////////////////Need not MOCK function ////////////////////////////////////////
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanesPgwNoConfig" );
        disconnect();
		reportTestResult("PostUserPlanesTest_PGWNOTCONFIG",
                           HTTP_SC_INVALID_UERPLANE_PROPERTISE, "OK",
                           status_code, resp_body["result"]);

        /////////////////////UserplaneAdd::API Check with wrong json request////////////////////////////////////////
        /////////////////////Need not MOCK function ////////////////////////////////////////
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanesSgwNoConfig" );
        disconnect();
		reportTestResult("PostUserPlanesTest_SGWNOTCONFIG",
                           HTTP_SC_INVALID_UERPLANE_PROPERTISE, "OK",
                           status_code, resp_body["result"]);		



        /////////////////////FcgiBackend::run return////////////////////////////////////////
        /////////////////////Need not MOCK function ////////////////////////////////////////
        // POST User plane ID 5 (URI invalid) 
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL_INVALIDURI,"", "PostUserplanes" );
        disconnect();
		reportTestResult("PostUserPlanesTest_INVALIDURI",
                           HTTP_SC_NOT_FOUND, "OK",
                           status_code, resp_body["result"]);



        // POST User plane ID 5 (URI invalid) 
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTBadRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanes" );
        disconnect();
		reportTestResult("PostUserPlanesTest_BADREQUEST",
                           HTTP_SC_BAD_REQUEST, "OK",
                           status_code, resp_body["result"]);
		

        // POST user plane ID  without s5u_pgw
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanesPgwNoS5U" );
        disconnect();
		reportTestResult("PostUserPlanesTest_PGWNOS5U",
                           HTTP_SC_INVALID_UERPLANE_PROPERTISE, "OK",
                           status_code, resp_body["result"]);


        // POST user plane ID  without s5u_pgw - up_ipaddress
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanesPgwNoS5uIpAddress" );
        disconnect();
		reportTestResult("PostUserPlanesTest_PGWNOS5U_IPADDRESS",
                           HTTP_SC_INVALID_UERPLANE_PROPERTISE, "OK",
                           status_code, resp_body["result"]);
		
        // POST user plane ID  without s5u_pgw - selector size = 0
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanesPgwZeroSelector" );
        disconnect();
		reportTestResult("PostUserPlanesTest_PGWZERO_SELECTOR",
                           HTTP_SC_INVALID_UERPLANE_PROPERTISE, "OK",
                           status_code, resp_body["result"]);


        // POST user plane ID  without selector
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanesPgwNoSelector" );
        disconnect();
		reportTestResult("PostUserPlanesTest_PGWNOSELECTOR",
                           HTTP_SC_INVALID_UERPLANE_PROPERTISE, "OK",
                           status_code, resp_body["result"]);		

        // POST user plane ID  without MNC
        status_code = "";
        resp_body.clear();		
        connect(TesterBase::host_ip_addr, TesterBase::host_port_num);

		sendPOSTRequest(status_code, resp_body, TesterBase::cookies["null"], POST_URL,"", "PostUserplanesPgwNoMNC" );
        disconnect();
		reportTestResult("PostUserPlanesTest_PGWNOMNC",
                           HTTP_SC_INVALID_UERPLANE_PROPERTISE, "OK",
                           status_code, resp_body["result"]);		


        return 0;
    } catch(Exception &e) {
        additionalMessage = e.err;
        disconnect();
        return -1;
    }
}
