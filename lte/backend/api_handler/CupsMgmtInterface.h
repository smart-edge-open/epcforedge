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

#define CUPS_API_INT64_TYPE 1
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

