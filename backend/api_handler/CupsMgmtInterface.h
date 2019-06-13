#ifndef __OAMAGENT__CUPSMGMTINTERFACE__
#define __OAMAGENT__CUPSMGMTINTERFACE__

#include <stdio.h>
#include <json/json.h>
#include <map>
#include <iostream>
#include <map>
#include <vector>
#include <boost/thread.hpp>
#include <json/json.h>

using namespace std;

/* ------------------------------------------------------------------------- */
/* Class Defs */
/* ------------------------------------------------------------------------- */

class CupsMgmtMessage 
{
public:

    vector<string> utilsSplitStr( string inputStr, char delimiter);
    int utilsParseApnNi(string apn_ni, string &apn, string &mnc, string &mcc);
	#ifdef CUPS_API_INT64_TYPE
	int utilsConvertStringToInt64(string value_hex);
	#endif
    int fillGetPgwResponse(Json::Value &pgwData, int pgwItemIndex,  Json::Value &response);
    int fillGetSgwResponse(Json::Value &sgwData, int sgwItemIndex,  Json::Value &response);	
	int fillPostPgwRequest(Json::Value &request, string &pgwPostData);
	int fillPostSgwRequest(Json::Value &request, string &sgwPostData);
    int fillPutPgwRequest(Json::Value &request, string &pgwPutData);
    int fillPutSgwRequest(Json::Value &request, string &sgwPutData);

};

#endif

