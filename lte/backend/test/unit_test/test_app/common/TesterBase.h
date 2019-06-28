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

//
// Created by david on 16-2-1.
// Modified by Yuan on 2-16-2017.
//

#ifndef MECFCGI_TESTERBASE_H
#define MECFCGI_TESTERBASE_H

/*******************************************************************************
* include declarations
*******************************************************************************/
#include <iostream>
#include <vector>
#include <map>

#include <boost/asio.hpp>
#include <boost/asio/ssl.hpp>
//#include <hiredis/hiredis.h>
#include <json/json.h>



/*******************************************************************************
* namespace expansion
*******************************************************************************/
using namespace std;
using namespace boost::asio;

/*******************************************************************************
* class definition/declaration
*******************************************************************************/

class TesterBase
{
    io_service svc;
    ssl::context ssl_ctx;
    ssl::stream<ip::tcp::socket> sock;

public:
    static const string host_ip_addr;
    static const int host_port_num;
    static const string verifCertPath;
    static const string hostName;

    static map<string,string> cookies;
    static map<string,string> provisions; // service name to service ID
    static map<string,string> traffic_rules; // restricted and unrestricted
    static map<string,string> subscriptions; //

public:
    typedef vector<string> ArrayType;

    // instantiate a TesterBase
    TesterBase();
    virtual ~TesterBase() {}

// execute test case, and receive execution status and additional message
// explaining the status code.
// purely virtual.
    virtual int execute (string &additionalMessage) = 0;

// establish a connection to a predefined target socket.
    void connect();
// establish a connection a remote target socket.
    int connect(string url, int port);
// close the established connection to the remote target scoket.
    int disconnect();

// send a POST request
    void post (const string &url, const string &data,  string _cookie="");
// send a POST request
    void postBad (const string &url, const string &data,  string _cookie="");

// send a POST request
    void patch (const string &url, const string &data,  string _cookie="");

// send a PUT request
    void put (const string &url, const string &data,  string _cookie="");
// send a GET request
    void get (const string &url, string _cookie="");
// send a DELETE request
    void del (const string &url, string _cookie="");

// record response
    int recordResponse (string &response);
// extract cookie header from response
    void extractCookie (string& cookie, const string& response);
// extract status code from response
    void extractStatusCode (string& statusCode, const string& response);
// extract response body, a string in JSON in our case, fromn response.
    void extractBody(string& body, const string& response);
// not sure what this does, but it is not used.
    void eventResponse (string &type, string &id);



    static bool addProvision(const string& serviceName, const string& serviceID);
    static bool addTrafficRule(const string& ruleName, const string& ruleID);
    static bool addSubscription(const string& subsName, const string& serviceID);
    static void printProvisions();
    static void printTrafficRules();
    static void printSubscriptions();
    static void printState();
    static void printCookie();

    /***************
     Convenience functions
    ***************/

// gets status code, response body and cookie by sending a POST request
// containg a body read from a .json file, and a custom cookie header
// to a URL.
    int sendPOSTRequest(
                string& statusCode, Json::Value& respBody, string& outCookie,
        const string& URL, const string& inCookie, const string& jsonFileName);
    int sendPOSTBadRequest(
                string& statusCode, Json::Value& respBody, string& outCookie,
        const string& URL, const string& inCookie, const string& jsonFileName);

   int sendPATCHRequest(
			string& statusCode, Json::Value& respBody, string& outCookie,
	  const string& URL, const string& inCookie, const string& jsonFileName);

// gets status code, response body and cookie by sending a PUT request
// containg a body read from a .json file, and a custom cookie header
// to a URL.
    int sendPUTRequest(
                string& statusCode, Json::Value& respBody, string& outCookie,
        const string& URL, const string& inCookie, const string& jsonFileName);

// gets status code, response body and cookie by sending a GET request
// containg a custom cookie header to a URL.
    int sendGETRequest(
                string& statusCode, Json::Value& respBody, string& outCookie,
                const string& URL, const string& inCookie);

// gets status code, response body and cookie by sending a DELETE request
// containg a custom cookie header to a URL.
    int sendDELRequest(
                string& statusCode, Json::Value& respBody, string& outCookie,
                const string& URL, const string& inCookie);

    int reportTestResult(const string& message,
                            const string& targetStatusCode,
                            const string& targetResult,
                            const string& actualStatusCode,
                            const Json::Value& actualResult);

};


#endif //MECFCGI_TESTERBASE_H
