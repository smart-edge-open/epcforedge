//
// Created by david on 16-2-1.
// Modified by Yuan on 02-14-2017
//

// library include
#include <boost/algorithm/string.hpp>

// project-wise include
#include "TesterBase.h"
#include "Exception.h"
#include "TestUtility.h"

using namespace boost::algorithm;
/*******************************************************************************
* static variable initializations
*******************************************************************************/
const string TesterBase::host_ip_addr {"127.0.0.1"};
const int TesterBase::host_port_num {8080};
const string TesterBase::verifCertPath {"/etc/nginx/ssl/mec.crt"};
const string TesterBase::hostName {"mec.local"};

map<string,string> TesterBase::cookies {}; // to be deprecated
map<string,string> TesterBase::provisions {};
map<string,string> TesterBase::traffic_rules {};
map<string,string> TesterBase::subscriptions {};

// constructor
TesterBase::TesterBase(): ssl_ctx(ssl::context::tlsv1), sock(svc, ssl_ctx) {
    boost::system::error_code ec;

    ssl_ctx.load_verify_file(verifCertPath, ec);
    if(ec) {
        std::cout << "in TesterBase::TesterBase - load_verify_file fails; error code: ";
        std::cout << boost::system::system_error(ec).what() << std::endl;
    }

    sock.set_verify_mode(boost::asio::ssl::verify_peer, ec);
    if(ec) {
        std::cout << "in TesterBase::TesterBase - set_verify_mode fails; error code: ";
        std::cout << boost::system::system_error(ec).what() << std::endl;
    }
}


// establish pre-defined socket connection
void TesterBase::connect() {
    // what is the prototype
    boost::system::error_code ec;
    sock.lowest_layer().connect({ {}, 8080 });
    if(ec) {
        std::cout << "in TesterBase::connect - connect fails; error code: ";
        std::cout << boost::system::system_error(ec).what() << std::endl;
    }
    sock.handshake(ssl::stream_base::handshake_type::client, ec);
    if(ec) {
        std::cout << "in TesterBase::connect - handshake fails; error code: ";
        std::cout << boost::system::system_error(ec).what() << std::endl;
    }
}


// establish a socket connection
int TesterBase::connect(string url, int port) {
    ip::tcp::endpoint ep(ip::address_v4::from_string(url), port);
    boost::system::error_code ec;
    sock.lowest_layer().connect(ep,ec);
    if(ec) {
        std::cout << "in TesterBase::connect - connect fails; error code: ";
        std::cout << boost::system::system_error(ec).what() << std::endl;
        return 1;
    }
    sock.handshake(ssl::stream_base::handshake_type::client, ec);
    if(ec) {
        std::cout << "in TesterBase::connect - handshake fails; error code: ";
        std::cout << boost::system::system_error(ec).what() << std::endl;
        return 1;
    }
    return 0;
}


int TesterBase::disconnect()
{
    boost::system::error_code ec;
    sock.shutdown(ec);
    if(ec) {
        // Ignoring short read error, not a real error
        if (ec.category() != boost::asio::error::get_ssl_category() ||
                ec.value() != ERR_PACK(ERR_LIB_SSL, 0, SSL_R_SHORT_READ)) {
            std::cout << "in TesterBase::disconnect - shutdown fails; error code: ";
            std::cout << boost::system::system_error(ec).what() << std::endl;
            sock.lowest_layer().close();
            return 1;
        }
    }
    sock.lowest_layer().close();
    return 0;
}


// send a POST request
void TesterBase::post(const string &url, const string &data, string _cookie) {
    stringstream ss;

    ss << "POST " << url << " HTTP/1.1\r\n";
    ss << "Accept: application/json\r\n";
    ss << "Content-Length: " << data.length() << "\r\n";
    ss << "Content-Type: application/json\r\n";
    if(_cookie.length()) {
        ss << "Cookie: " << _cookie << "\r\n";
    }
    ss << "Host: " << hostName << "\r\n\r\n";
    ss << data;

    boost::asio::write(sock, buffer(ss.str()));
}
// send a POST request
void TesterBase::postBad(const string &url, const string &data, string _cookie) {
    stringstream ss;

    ss << "POSTBAD " << url << " HTTP/1.1\r\n";
    ss << "Accept: application/json\r\n";
    ss << "Content-Length: " << data.length() << "\r\n";
    ss << "Content-Type: application/json\r\n";
    if(_cookie.length()) {
        ss << "Cookie: " << _cookie << "\r\n";
    }
    ss << "Host: " << hostName << "\r\n\r\n";
    ss << data;

    boost::asio::write(sock, buffer(ss.str()));
}

// send a PATCH request
void TesterBase::patch(const string &url, const string &data, string _cookie) {
    stringstream ss;

    ss << "PATCH " << url << " HTTP/1.1\r\n";
    ss << "Accept: application/json\r\n";
    ss << "Content-Length: " << data.length() << "\r\n";
    ss << "Content-Type: application/json\r\n";
    if(_cookie.length()) {
        ss << "Cookie: " << _cookie << "\r\n";
    }
    ss << "Host: " << hostName << "\r\n\r\n";
    ss << data;

    boost::asio::write(sock, buffer(ss.str()));
}


// send a PUT request
void TesterBase::put(const string &url, const string &data, string _cookie) {
    stringstream ss;

    ss << "PUT " << url << " HTTP/1.1\r\n";
    ss << "Accept: application/json\r\n";
    ss << "Content-Length: " << data.length() << "\r\n";
    ss << "Content-Type: application/json\r\n";
    if(_cookie.length()) {
        ss << "Cookie: " << _cookie << "\r\n";
    }
    ss << "Host: " << hostName << "\r\n\r\n";
    ss << data;
    boost::asio::write(sock, buffer(ss.str()));
}


// send a GET request
void TesterBase::get(const string &url, string _cookie) {
    stringstream ss;
    ss << "GET " << url << " HTTP/1.1\r\n";
    if(_cookie.length()) {
        ss << "Cookie: " << _cookie << "\r\n";
    }
    ss << "Host: " << hostName << "\r\n\r\n";
    boost::asio::write(sock, buffer(ss.str()));
}


// send a DELETE request, cannot use delete for func name because it is keyword
void TesterBase::del(const string &url, string _cookie) {
    stringstream ss;
    ss << "DELETE " << url << " HTTP/1.1\r\n";
    if(_cookie.length()) {
        ss << "Cookie: " << _cookie << "\r\n";
    }
    ss << "Host: " << hostName << "\r\n\r\n";
    boost::asio::write(sock, buffer(ss.str()));
}


// record response from socket, returns exit code 1 on error, 0 on normal.
int TesterBase::recordResponse(string &response) {
    try {
            response = "";
            boost::system::error_code ec;
            while(true) {
                array<char, 81920> recvBuf;
                size_t data_read = sock.read_some(buffer(recvBuf), ec);
                if (ec) {
                    std::cout << "in TesterBase::recordResponse - sock.read_some fails; error code: ";
                    std::cout << boost::system::system_error(ec).what() << std::endl;
                    break;
                }
                response.append(recvBuf.data(), data_read);
                if(response.find("Transfer-Encoding: chunked") == string::npos
                    || response.find("\r\n0\r\n\r\n") != string::npos) {
                    break;
                }
            }
            return 0;
    }
    catch(boost::system::system_error &e){
        std::cout << "in TesterBase::recordResponse: exception is thrown- ";
        std::cout << e.what() << std::endl;
        int posStart, posEnd;
        string statusCode = "";
        if( (posStart = response.find("HTTP/1.1 ")) != string::npos) {
            if( (posEnd = response.find("\r\n", posStart)) != string::npos) {
                statusCode = response.substr(posStart + 9, posEnd - posStart - 9);
                trim(statusCode);
            }
        }
        if (0 != statusCode.compare("200 OK")) {
            response = "";
            return 1;
        } else {
            return 0;
        }
    }
}


// extract cookie from recorded response
void TesterBase::extractCookie(string& cookie, const string& response) {
    int posStart, posEnd;
    if( (posStart = response.find("Set-Cookie:")) != string::npos) {
        if( (posEnd = response.find("\r\n", posStart)) != string::npos) {
                cookie = response.substr(posStart + 11, posEnd - posStart - 11);
                trim(cookie);
        }
    }
}


// extract status code from recorded response
void TesterBase::extractStatusCode(string& statusCode, const string& response) {
    int posStart, posEnd;
    if( (posStart = response.find("HTTP/1.1 ")) != string::npos) {
        if( (posEnd = response.find("\r\n", posStart)) != string::npos) {
            statusCode = response.substr(posStart + 9, posEnd - posStart - 9);
            trim(statusCode);
        }
    }
}


// extract body from recorded response
void TesterBase::extractBody(string& body, const string& response) {
    int posStart = response.find("{");
    if(posStart != string::npos) {
        int posEnd = response.rfind("}");
        if(posEnd == string::npos) {
            throw Exception(Exception::PARSING_JSON_BODY);
        }
        body = response.substr(posStart, posEnd - posStart + 1);
    }
}


bool TesterBase::addProvision(const string& serviceName, const string& serviceID) {
    return (TesterBase::provisions.emplace(serviceName, serviceID)).second;
}


bool TesterBase::addTrafficRule(const string& ruleName, const string& ruleID) {
    return (TesterBase::traffic_rules.emplace(ruleName, ruleID)).second;
}


bool TesterBase::addSubscription(const string& subsName, const string& serviceID) {
    return (TesterBase::subscriptions.emplace(subsName, serviceID)).second;
}


void TesterBase::printProvisions() {
    std::cout << " -- TesterBase::provisions := [ \n";
    for (auto & provision : TesterBase::provisions) {

        std::cout
        << " ---- "
        << provision.first << " <-> " << provision.second
        << "\n";
    }
    std::cout << " -- ]"<< std::endl;
}


void TesterBase::printTrafficRules() {
    std::cout << " -- TesterBase::traffic_rules := [ \n";
    for (auto & traffic_rule : TesterBase::traffic_rules) {

        std::cout
        << " ---- "
        << traffic_rule.first << " <-> " << traffic_rule.second
        << "\n";
    }
    std::cout << " -- ]"<< std::endl;
}


void TesterBase::printSubscriptions() {
    std::cout << " -- TesterBase::subscriptions := [ \n";
    for (auto& subscription : TesterBase::subscriptions) {

        std::cout
        << " ---- "
        << subscription.first << " <-> " << subscription.second
        << "\n";
    }
    std::cout << " -- ]"<< std::endl;
}


void TesterBase::printState() {
    std::cout
    << " - TestBase internal state:\n";
    TesterBase::printCookie();
    TesterBase::printProvisions();
    TesterBase::printSubscriptions();
    TesterBase::printTrafficRules();
}


void TesterBase::printCookie() {
    std::cout << " -- TesterBase::cookies := [ \n";
    for (auto& cookie : TesterBase::cookies) {

        std::cout
        << " ---- "
        << cookie.first << " <-> " << cookie.second
        << "\n";
    }
}


/*******************************************************************************
* convenience functions
*******************************************************************************/

// gets status code, response body and cookie by sending a POST request
// containg a body read from a .json file to a URL.
int TesterBase::sendPOSTRequest(
        string& statusCode, Json::Value& respBody, string& outCookie,
        const string& URL, const string& inCookie, const string& jsonFileName) {

    string req_body, response, resp_body;
    int return_code;
    Json::Reader reader;
    req_body = JSONFileToString(jsonFileName);
	post(URL, req_body, inCookie);
	return_code = recordResponse(response);
    if(0 != return_code) {
	    return 1;
    }
	//printf("extracting cookie\n");		
    //extractCookie(outCookie, response);
	extractStatusCode(statusCode, response);
	extractBody(resp_body, response);
	
    reader.parse(resp_body, respBody);

    return 0;
}

int TesterBase::sendPOSTBadRequest(
        string& statusCode, Json::Value& respBody, string& outCookie,
        const string& URL, const string& inCookie, const string& jsonFileName) {

    string req_body, response, resp_body;
    int return_code;
    Json::Reader reader;
    req_body = JSONFileToString(jsonFileName);
	postBad(URL, req_body, inCookie);
	return_code = recordResponse(response);
    if(0 != return_code) {
	    return 1;
    }
	//printf("extracting cookie\n");		
    //extractCookie(outCookie, response);
	extractStatusCode(statusCode, response);
	extractBody(resp_body, response);
	
    reader.parse(resp_body, respBody);

    return 0;
}

int TesterBase::sendPATCHRequest(
        string& statusCode, Json::Value& respBody, string& outCookie,
        const string& URL, const string& inCookie, const string& jsonFileName) {

    string req_body, response, resp_body;
    int return_code;
    Json::Reader reader;
    req_body = JSONFileToString(jsonFileName);
	patch(URL, req_body, inCookie);
	return_code = recordResponse(response);
    if(0 != return_code) {
	    return 1;
    }
	//printf("extracting cookie\n");		
    //extractCookie(outCookie, response);
	extractStatusCode(statusCode, response);
	extractBody(resp_body, response);
	
    reader.parse(resp_body, respBody);

    return 0;
}

// gets status code, response body and cookie by sending a PUT request
// containg a body read from a .json file to a URL.
int TesterBase::sendPUTRequest(
        string& statusCode, Json::Value& respBody, string& outCookie,
        const string& URL, const string& inCookie, const string& jsonFileName) {

    string req_body, response, resp_body;
    int return_code;
    Json::Reader reader;

    req_body = JSONFileToString(jsonFileName);
    put(URL, req_body, inCookie);
    return_code = recordResponse(response);
    if(0 != return_code) {
        return 1;
    }
   // extractCookie(outCookie, response);
    extractStatusCode(statusCode, response);
    extractBody(resp_body, response);

    reader.parse(resp_body, respBody);

    return 0;
}


// gets status code, response body and cookie by sending a GET request to a
// URL.
int TesterBase::sendGETRequest(
                string& statusCode, Json::Value& respBody, string& outCookie,
                const string& URL, const string& inCookie) {

    string response, resp_body;
    int return_code;
    Json::Reader reader;

    get(URL, inCookie);
    return_code = recordResponse(response);
    if(0 != return_code) {
        return 1;
    }
    //extractCookie(outCookie, response);
    extractStatusCode(statusCode, response);
    extractBody(resp_body, response);

    reader.parse(resp_body, respBody);

    return 0;
}


// gets status code, response body and cookie by sending a DELETE request to
// a URL.
int TesterBase::sendDELRequest(
                string& statusCode, Json::Value& respBody, string& outCookie,
                    const string& URL, const string& inCookie) {

    string response, resp_body;
    int return_code;
    Json::Reader reader;
    del(URL, inCookie);
    return_code = recordResponse(response);
    if(0 != return_code) {
		printf("report c\n");
        return 1;
    }
    //extractCookie(outCookie, response);
    extractStatusCode(statusCode, response);
    extractBody(resp_body, response);
    reader.parse(resp_body, respBody);
    return 0;
}

int okNum = 0;
int ngNum = 0;

int TesterBase::reportTestResult(const string& message,
                            const string& targetStatusCode,
                            const string& targetResult,
                            const string& actualStatusCode,
                            const Json::Value& actualResult) {
    // Disable result compare
    //if(0 == actualStatusCode.compare(targetStatusCode) &&
    //        !actualResult.empty() &&
    //    0 == actualResult.asString().compare(targetResult))

    if(0 == actualStatusCode.compare(targetStatusCode))	{
        cout << "[        OK] " << message << ": (" << actualStatusCode  << ")" << endl;
        okNum++;
		return 1;
    } else {
        cout << "[        NG] " << message << ": (" << actualStatusCode  << ")" << endl;
        ngNum++;
		return 0;
    }
}
