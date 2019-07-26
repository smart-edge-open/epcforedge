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
/**
 * @file    cupsMgmtInterface.cpp
 * @brief   Implementation of CUPS management interface between MEC 
 *          controller and EPC OAMAgent
 ********************************************************************/

#include <iostream>
#include <map>
#include <vector>
#include <boost/thread.hpp>
#include <json/json.h>


#include "UpfController.h"
#include "Exception.h"
#include "HandlerCommon.h"
#include "Log.h"
#include "RawRequest.h"
#include "LocalConfig.h"
#include "CpfInterface.h"
#include "CupsMgmtInterface.h"



/**
* @brief       The utility function will split string according to delimiter  
*					
* @param[in]   inputStr      The string to be splitted.
* @param[in]   delimiter     The delimiter for splitting.
* @return      String Vector with splitted strings.
*/

vector<string> CupsMgmtMessage::utilsSplitStr( string inputStr, char delimiter)
{
   vector <string> outputVect;
   stringstream sStream(inputStr);
   string val;
   while(getline(sStream, val, delimiter)) {
      outputVect.push_back(val);
   }
   return outputVect;
}

/**
* @brief			The utility function will parse apnni element from EPC  
*					
* @param[in]		apn_ni 	  APNNI Data from EPC
* @param[out]		apn       Parsed APN value.
* @param[out]		mnc       Parsed mnc value.
* @param[out]		mcc       Parsed mcc value.
* @return			0:success; 1:failure
*/
int CupsMgmtMessage::utilsParseApnNi(string apn_ni, string &apn, string &mnc, string &mcc)
{

    vector<string> apn_ni_vect = utilsSplitStr(apn_ni, '.');
    int size = apn_ni_vect.size();
    if (4 != size) {
       OAMAGENT_LOG(ERR, "not valid size %d.\n", size);
       return -1;
    }
    OAMAGENT_LOG(INFO, "Get apn %s, mnc %s, mcc %s, domain %s.\n", 
		apn_ni_vect[0].c_str(), apn_ni_vect[1].c_str(),
		apn_ni_vect[2].c_str(), apn_ni_vect[3].c_str());
    apn = apn_ni_vect[0];
    mnc = apn_ni_vect[1].substr(3);
    mcc = apn_ni_vect[2].substr(3);	
    return 0;
}

#ifdef CUPS_API_INT64_TYPE
/**
* @brief        The utility function will convert string with hex format to int64 
*					
* @param[in]	value_hex 	The string with value as HEX format
* @return	int
*/

int  CupsMgmtMessage::utilsConvertHexStringToDecInt(string value_hex)
{
    int res = std::stoi (value_hex,0,16);
    OAMAGENT_LOG(INFO, "Convert (%s) to (%d) \n", value_hex.c_str(), res);
    return res;
}
/**
* @brief        The utility function will convert decimal int to hex string 
*					
* @param[in]	value_dec 	The decimal value
* @return	hex base string
*/

std::string  CupsMgmtMessage::utilsConvertDecIntToHexString(int value_dec)
{
    std::string ret;
	std::stringstream temp;
	temp << std::hex << value_dec;
    temp >> ret;
    return ret;
}

#endif

/**
* @brief         The inteface function will fill PGW Infor for GET Response to MECController  
*					
* @param[in]     pgwData       PGW List Data from EPC
* @param[in]     pgwItemIndex  PGW List Index.
* @param[out]    response      JSON-formatted key-value pair(s) for GETResponse to MECController
* @return        0:success; 1:failure
*/

int CupsMgmtMessage::fillGetPgwResponse(Json::Value &pgwData, int pgwItemIndex, Json::Value &response)
{

    OAMAGENT_LOG(INFO,"Checking Input PGW JSON Data.\n");

    // check GET body from EPC
    if(pgwData["items"][pgwItemIndex]["s5u_ip"].isString()){
    }
    else {		
    	OAMAGENT_LOG(ERR, "s5u_ip not exist.\n");	
        return -1;
    }		
    if(pgwData["items"][pgwItemIndex]["tac"].isString()){
    }
    else {
    	OAMAGENT_LOG(ERR, "tac not exist.\n");	
        return -1;
	}
    if(pgwData["items"][pgwItemIndex]["apn_ni"].isString()){
    }
    else {
    	OAMAGENT_LOG(ERR, "apn_ni not exist.\n");	
        return -1;
    }
    if(pgwData["items"][pgwItemIndex]["uuid"].isString()){
    }
    else {
    	OAMAGENT_LOG(ERR, "uuid not exist.\n");	
        return -1;
    }
	
    response["uuid"] = pgwData["items"][pgwItemIndex]["uuid"];
    response["id"] = pgwData["items"][pgwItemIndex]["id"];
    response["function"] = "SAEGWU";

     //STEP1: fill config
    Json::Value config;
    config["s5u_pgw"]["up_ip_address"] = pgwData["items"][pgwItemIndex]["s5u_ip"];		
    
    //STEP2: fill selectors
    Json::Value selectors_array;    
    Json::Value selectors;
    Json::Value network;
	//STEP2.0 selectors->id
	selectors["id"] = pgwData["items"][pgwItemIndex]["id"];	
	
	//STEP2.1 selectors->network
    string apn_ni,apn,mnc,mcc;
    apn_ni = pgwData["items"][pgwItemIndex]["apn_ni"].asString();
    if (0 != utilsParseApnNi( apn_ni,apn,mnc,mcc)) {
       return -1;
    }
    network["mcc"] = mcc;
    network["mnc"] = mnc;		
    selectors["network"] = network; 	
	//STEP2.1 selectors->uli	
    Json::Value uli;
    Json::Value tai;
	#ifdef CUPS_API_INT64_TYPE
    // convert into int64 type (JSON not support HEX, so need to convert into DEC)
    // BUT TAC only has 16 bit. so int type will be enough
    tai["tac"] = utilsConvertHexStringToDecInt(pgwData["items"][pgwItemIndex]["tac"].asString());
	#else
	// default should be string type.
    tai["tac"] = pgwData["items"][pgwItemIndex]["tac"];
	#endif
    uli["tai"] = tai;
	selectors["uli"] = uli;
	//STEP2.3 selectors->pdn
    Json::Value pdn;	 
    pdn["apns"].append(apn.c_str());	
    selectors["pdn"] = pdn;		

    // Finally, put all into reponse
    response["config"] = config;		
    selectors_array.append(selectors);	
    response["selectors"] = selectors_array;

	return 0;
     
}

/**
* @brief        The inteface function will fill SGW Infor for GET Response to MECController  
*					
* @param[in]    sgwData       SGW List Data from EPC
* @param[in]    sgwItemIndex  SGW List Index.
* @param[out]   response      JSON-formatted key-value pair(s) for GETResponse to MECController
* @return       0:success; 1:failure
*/

int CupsMgmtMessage::fillGetSgwResponse(Json::Value &sgwData, int sgwItemIndex, Json::Value &response)
{

    OAMAGENT_LOG(INFO,"Checking Input SGW JSON Data.\n");

    if(sgwData["items"][sgwItemIndex]["s5u_ip"].isString()){
    }
    else {		
    	OAMAGENT_LOG(ERR, "s5u_ip not exist.\n");	
        return -1;
    }		
    if(sgwData["items"][sgwItemIndex]["tac"].isString()){
    }
    else {
    	OAMAGENT_LOG(ERR, "tac not exist.\n");	
        return -1;
	}
    if(sgwData["items"][sgwItemIndex]["s1u_ip"].isString()){
    }
    else {
    	OAMAGENT_LOG(ERR, "s1u_ip not exist.\n");	
        return -1;
    }
    if(sgwData["items"][sgwItemIndex]["uuid"].isString()){
    }
    else {
    	OAMAGENT_LOG(ERR, "uuid not exist.\n");	
        return -1;
    }
	
    response["uuid"] = sgwData["items"][sgwItemIndex]["uuid"];
    response["id"] = sgwData["items"][sgwItemIndex]["id"];
    //STEP1: fill config
    //Json::Value config;
    response["config"]["s5u_sgw"]["up_ip_address"] = sgwData["items"][sgwItemIndex]["s5u_ip"];
    response["config"]["s1u"]["up_ip_address"]     = sgwData["items"][sgwItemIndex]["s1u_ip"];	
    return 0;
     
}

/**
* @brief        The inteface function will fill PGW Infor for POST Request to EPC 
*					
* @param[in]    request 	  JSON-formatted key-value pair(s), PostRequest from MECController
* @param[out]   pgwPostData	  JSON-formatted key-value pair(s) for POSTRequest to EPC
* @return       0:success; 1:failure
*/

int CupsMgmtMessage::fillPostPgwRequest(Json::Value &request, string &pgwPostData)
{
   int size;
   Json::Value pgwPostJson;
   //stringstream jsonStr;
   string userplaneID = request.get("id", "Nil").asString();
   if (0 == userplaneID.compare("Nil")) {
      OAMAGENT_LOG(WARNING, "[id] is not found in request.\n");
   }
   //OAMAGENT_LOG(INFO, "UserplaneAdd userplaneID (%s).\n", userplaneID.c_str());

   // get uuid and put into POST to EPC
   Json::Value uuid = request.get("uuid","Nil");
   if (0 == uuid.compare("Nil")) {
       OAMAGENT_LOG(ERR, "[uuid] is not found in request.\n");   	
       return -1;
   }
   pgwPostJson["uuid"] = uuid.asString();
   

   /////////////////start read config configuration ////////////////////////////////////////		
   Json::Value config = request.get("config","Nil");
   if (0 == config.compare("Nil")) {
       return -1;
   }

   // get s5u_pgw configuration
   Json::Value s5u_pgwCfg = config.get ("s5u_pgw","Nil");
   if (0 == s5u_pgwCfg.compare("Nil")) {
       OAMAGENT_LOG(ERR, "[s5u_pgwCfg] is not found in request.\n");
       return -1;
   }		
   string s5uPgwIp_up  = s5u_pgwCfg.get ("up_ip_address","Nil").asString();			
   if (0 == s5uPgwIp_up.compare("Nil")) {
       OAMAGENT_LOG(ERR, "[s5uPgwIp_up] is not found in request.\n");
       return -1;
   }
   
   pgwPostJson["s5u_ip"] = s5uPgwIp_up;
   pgwPostJson["proxy_ip"] = s5uPgwIp_up;
   OAMAGENT_LOG(INFO, "UserplaneAdd s5u_pgw up_ip_address (%s).\n", s5uPgwIp_up.c_str());

   /////////////////start read selectors configuration ////////////////////////////////////////
   Json::Value selectors = request.get("selectors","Nil");
   if (0 == selectors.compare("Nil")) {
       OAMAGENT_LOG(ERR, "[selectors] is not found in request.\n");
       return -1;
   }
   size = selectors.size();
   if (0 == size) {
       OAMAGENT_LOG(ERR, "[selectors] is not found.\n");
       return -1;
   }
   OAMAGENT_LOG(INFO, "UserplaneAdd selectors size (%d).\n", size);
   
   for (int i = 0; i < size; i++){
       string selectors_id = selectors[i].get ("id","Nil").asString();
       if (0 == selectors_id.compare("Nil")){
           OAMAGENT_LOG(WARNING, "[selectors_id] is not found in breakout.\n");
           //return -1;
       }
       else {
           OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] id (%s).\n", i, selectors_id.c_str());
       }
       /////// network
       Json::Value network = selectors[i].get ("network","Nil");
       string network_mcc = network.get ("mcc","Nil").asString();
       string network_mnc = network.get ("mnc","Nil").asString();
       if (0 == network_mnc.compare("Nil") || 0 == network_mcc.compare("Nil")) {
           OAMAGENT_LOG(ERR, "[network] is not found in breakout.\n");
           return -1;
       }
       OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] network (mcc-%s, mnc-%s).\n", i, network_mcc.c_str(), network_mnc.c_str());
   
       /////// uli
       Json::Value uli = selectors[i].get ("uli","Nil");
       Json::Value tai = uli.get ("tai","Nil");
       #ifdef CUPS_API_INT64_TYPE
       int tac = tai.get ("tac","Nil").asInt();
       pgwPostJson["tac"] = utilsConvertDecIntToHexString(tac);
       OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] uli (tac %d).\n", i, tac);
       #else
       string tac = tai.get ("tac","Nil").asString();
       pgwPostJson["tac"] = tac;
       OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] uli (tac %s).\n", i, tac.c_str());
       #endif
       /////// pdn
       Json::Value pdn = selectors[i].get ("pdn","Nil");
       Json::Value apns = pdn.get ("apns","Nil");
       string apn_name;
       int ss = apns.size();
       OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] apn size (%d).\n",i, ss);
       if (ss <=0) {
           OAMAGENT_LOG(ERR, "not found apns setting.\n");
           return -1;
       }			
       for (int j = 0; j < ss; j++) {
           apn_name = apns[j].asString();
           OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] apn[%d] name (%s).\n",i,j, apn_name.c_str());				 
       }
       string apn_ni = apn_name + ".mnc" + network_mnc + ".mcc" + network_mcc + ".gprs";
       pgwPostJson["apn_ni"] = apn_ni;
   	
   }		

   OAMAGENT_LOG(INFO, "PGW_ADD JSON: %s\n",pgwPostJson.toStyledString().c_str());
   pgwPostData = pgwPostJson.toStyledString();

   return 0;
}

/**
* @brief        The inteface function will fill SGW Infor for POST Request to EPC 
*					
* @param[in]    request 	  JSON-formatted key-value pair(s), POSTRequest from MECController
* @param[out]   sgwPostData	  JSON-formatted key-value pair(s) for POSTRequest to EPC
* @return       0:success; 1:failure
*/

int CupsMgmtMessage::fillPostSgwRequest(Json::Value &request, string &sgwPostData)
{
    int size;
    Json::Value sgwPostJson;


    string userplaneID = request.get("id", "Nil").asString();
    if (0 == userplaneID.compare("Nil")) {
        OAMAGENT_LOG(WARNING, "[id] is not found in request.\n");
    }
    //OAMAGENT_LOG(INFO, "UserplaneAdd userplaneID (%s).\n", userplaneID.c_str());
	
    // get uuid and put into POST to EPC
    Json::Value uuid = request.get("uuid","Nil");
    if (0 == uuid.compare("Nil")) {
        OAMAGENT_LOG(ERR, "[uuid] is not found in request.\n");	
        return -1;
    }
    sgwPostJson["uuid"] = uuid.asString();		   

    /////////////////start read config configuration ////////////////////////////////////////		
    Json::Value config = request.get("config","Nil");
    if (0 == config.compare("Nil")) {
        OAMAGENT_LOG(ERR, "[config] is not found in request.\n");
        return -1;
    }
	
    // get s5u_pgw configuration
    Json::Value s5u_sgwCfg = config.get ("s5u_sgw","Nil");
    if (0 == s5u_sgwCfg.compare("Nil")) {
        OAMAGENT_LOG(ERR, "[s5u_SgwCfg] is not found in request.\n");
        return -1;
    }		
    string s5uSgwIp_up  = s5u_sgwCfg.get ("up_ip_address","Nil").asString();			
    if (0 == s5uSgwIp_up.compare("Nil")) {
        OAMAGENT_LOG(ERR, "[s5uSgwIp_up] is not found in request.\n");
        return -1;
    }
    sgwPostJson["s5u_ip"]   = s5uSgwIp_up;
    sgwPostJson["peer_ip"]   = s5uSgwIp_up;	
    sgwPostJson["s1u_nat_ip"] = s5uSgwIp_up;

    // get s1u configuration
    Json::Value s1uCfg = config.get ("s1u","Nil");
    if (0 == s1uCfg.compare("Nil")) {
        OAMAGENT_LOG(ERR, "[s1uCfg] is not found in request.\n");
        return -1;
    }		
    string s1uIp_up  = s1uCfg.get ("up_ip_address","Nil").asString();			
    if (0 == s1uIp_up.compare("Nil")) {
        OAMAGENT_LOG(ERR, "[s1uIp_up] is not found in request.\n");
        return -1;
    }
    sgwPostJson["s1u_ip"]   = s1uIp_up;
    OAMAGENT_LOG(INFO, "UserplaneAdd s5u_ip (%s), s1u_ip (%s).\n", s5uSgwIp_up.c_str(), s1uIp_up.c_str());

    /////////////////start read selectors configuration ////////////////////////////////////////
    Json::Value selectors = request.get("selectors","Nil");
    if (0 == selectors.compare("Nil")) {
        OAMAGENT_LOG(ERR, "[selectors] is not found in request.\n");
        return -1;
    }
    size = selectors.size();
    if (0 == size) {
        OAMAGENT_LOG(ERR, "[selectors] is not found.\n");
        return -1;
    }
    OAMAGENT_LOG(INFO, "UserplaneAdd selectors size (%d).\n", size);

    for (int i = 0; i < size; i++){
        string selectors_id = selectors[i].get ("id","Nil").asString();
        if (0 == selectors_id.compare("Nil")){
           OAMAGENT_LOG(WARNING, "[selectors_id] is not found in breakout.\n");
           //return -1;
        }
        else {
           OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] id (%s).\n", i, selectors_id.c_str());
        }

        /////// uli
        Json::Value uli = selectors[i].get ("uli","Nil");
        Json::Value tai = uli.get ("tai","Nil");
#ifdef CUPS_API_INT64_TYPE
        int tac = tai.get ("tac","Nil").asInt();
        sgwPostJson["tac"] = utilsConvertDecIntToHexString(tac);
        OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] uli (tac %d).\n", i, tac);
#else
        string tac = tai.get ("tac","Nil").asString();
        sgwPostJson["tac"] = tac;
        OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] uli (tac %s).\n", i, tac.c_str());

#endif
    }		
		 
    OAMAGENT_LOG(INFO, "SGW_ADD JSON: %s\n",sgwPostJson.toStyledString().c_str());
    sgwPostData = sgwPostJson.toStyledString();
    return 0;
}

/**
* @brief           The inteface function will fill PGW Infor for PUT Request to EPC 
*					
* @param[in]       request 	  JSON-formatted key-value pair(s), PatchRequest from MECController
* @param[out]      pgwPostData	  JSON-formatted key-value pair(s) for PUTRequest to EPC
* @return          0:success; 1:failure
*/

int CupsMgmtMessage::fillPutPgwRequest(Json::Value &request, string &pgwPutData)
{
   int size;
   Json::Value pgwPutJson;
   //stringstream jsonStr;
   string userplaneID = request.get("id", "Nil").asString();
   if (0 == userplaneID.compare("Nil")) {
       OAMAGENT_LOG(WARNING, "[id] is not found in request.\n");
   }
   //OAMAGENT_LOG(INFO, "UserplanePatch userplaneID (%s).\n", userplaneID.c_str());

   // get uuid and put into POST to EPC
   Json::Value uuid = request.get("uuid","Nil");
   if (0 == uuid.compare("Nil")) {
       OAMAGENT_LOG(ERR, "[uuid] is not found in request.\n");   	
       return -1;
   }
   pgwPutJson["uuid"] = uuid.asString();

   /////////////////start read config configuration ////////////////////////////////////////		
   Json::Value config = request.get("config","Nil");
   if (0 == config.compare("Nil")) {
       OAMAGENT_LOG(INFO, "UserplanePatch: not config element \n");
   }
   else {
       // get s5u_pgw configuration
       Json::Value s5u_pgwCfg = config.get ("s5u_pgw","Nil");
       if (0 == s5u_pgwCfg.compare("Nil")) {
           OAMAGENT_LOG(INFO, "[s5u_pgwCfg] is not found in request.\n");	
       }
       else {
           string s5uPgwIp_up  = s5u_pgwCfg.get ("up_ip_address","Nil").asString();			
           if (0 == s5uPgwIp_up.compare("Nil")) {
               OAMAGENT_LOG(ERR, "[s5uPgwIp_up] is not found in request.\n");
           }
           else {
               pgwPutJson["s5u_ip"]   = s5uPgwIp_up;
               pgwPutJson["proxy_ip"] = s5uPgwIp_up;
               OAMAGENT_LOG(INFO, "UserplaneAdd s5u_pgw up_ip_address (%s).\n", s5uPgwIp_up.c_str());
           }
       }
   }

   /////////////////start read selectors configuration ////////////////////////////////////////
   Json::Value selectors = request.get("selectors","Nil");
   if (0 == selectors.compare("Nil")) {
       OAMAGENT_LOG(INFO, "[selectors] is not found in request.\n");
   }
   else {
       size = selectors.size();
       OAMAGENT_LOG(INFO, "UserplaneAdd selectors size (%d).\n", size);	   
       if (0 < size) {
           for (int i = 0; i < size; i++){
              string selectors_id = selectors[i].get ("id","Nil").asString();
              /////// uli
              Json::Value uli = selectors[i].get ("uli","Nil");
              Json::Value tai = uli.get ("tai","Nil");
#ifdef CUPS_API_INT64_TYPE			  
              int tac = tai.get ("tac","Nil").asInt();
              pgwPutJson["tac"] = utilsConvertDecIntToHexString(tac);
              OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] uli (tac %d).\n", i, tac);
#else
              string tac = tai.get ("tac","Nil").asString();
              pgwPutJson["tac"] = tac;
              OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] uli (tac %s).\n", i, tac.c_str());
#endif
              /////// network
              Json::Value network = selectors[i].get ("network","Nil");
              string network_mcc = network.get ("mcc","Nil").asString();
              string network_mnc = network.get ("mnc","Nil").asString();
              if (0 == network_mnc.compare("Nil") || 0 == network_mcc.compare("Nil")) {
                 OAMAGENT_LOG(ERR, "[network] is not found in breakout.\n");
                 break;
              }
              OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] network (mcc-%s, mnc-%s).\n", i, network_mcc.c_str(), network_mnc.c_str());
			   
              /////// pdn
              Json::Value pdn = selectors[i].get ("pdn","Nil");
              Json::Value apns = pdn.get ("apns","Nil");
              string apn_name;
              int ss = apns.size();
              OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] apn size (%d).\n",i, ss);
              if (ss <=0) {
                 OAMAGENT_LOG(ERR, "not found apns setting.\n");
                 break;
              }	
              for (int j = 0; j < ss; j++) {
                 apn_name = apns[j].asString();
                 OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] apn[%d] name (%s).\n",i,j, apn_name.c_str());
              }
              string apn_ni = apn_name + ".mnc" + network_mnc + ".mcc" + network_mcc + ".gprs";
              pgwPutJson["apn_ni"] = apn_ni;			   
           }
        }		
    }
    //Json::Value jsonData;
    Json::Reader jsonReader;
    Json::FastWriter fastWriter;
    pgwPutData = fastWriter.write(pgwPutJson);

    OAMAGENT_LOG(INFO, "PGW_PUT JSON: %s\n",pgwPutJson.toStyledString().c_str());
    return 0;
}

/**
* @brief          The inteface function will fill SGW Infor for PUT Request to EPC 
*					
* @param[in]      request 	  JSON-formatted key-value pair(s), PatchRequest from MECController
* @param[out]     sgwPutData	  JSON-formatted key-value pair(s) for PUTRequest to EPC
* @return         0:success; 1:failure
*/

int CupsMgmtMessage::fillPutSgwRequest(Json::Value &request, string &sgwPutData)
{
    int size;
    Json::Value sgwPutJson;
    //stringstream jsonStr;		
    string userplaneID = request.get("id", "Nil").asString();
    if (0 == userplaneID.compare("Nil")) {
       OAMAGENT_LOG(WARNING, "[id] is not found in request.\n");
       //return -1;
    }
    //OAMAGENT_LOG(INFO, "UserplanePatch userplaneID (%s).\n", userplaneID.c_str());
    // get uuid and put into POST to EPC
    Json::Value uuid = request.get("uuid","Nil");
    if (0 == uuid.compare("Nil")) {
       OAMAGENT_LOG(ERR, "[uuid] is not found in request.\n");	 
       return -1;
    }
    sgwPutJson["uuid"] = uuid.asString();

    /////////////////start read config configuration ////////////////////////////////////////		
    Json::Value config = request.get("config","Nil");
    if (0 == config.compare("Nil")) {
       OAMAGENT_LOG(ERR, "[config] is not found in request.\n");			
    }
    else {
       // get s5u_pgw configuration
       Json::Value s5u_sgwCfg = config.get ("s5u_sgw","Nil");
       if (0 == s5u_sgwCfg.compare("Nil")) {
           OAMAGENT_LOG(ERR, "[s5u_SgwCfg] is not found in request.\n");
       }		
       else {
           string s5uSgwIp_up  = s5u_sgwCfg.get ("up_ip_address","Nil").asString();			
           if (0 == s5uSgwIp_up.compare("Nil")) {
               OAMAGENT_LOG(ERR, "[s5uSgwIp_up] is not found in request.\n");
           }
           else {
               sgwPutJson["s5u_ip"]	= s5uSgwIp_up;
               //sgwPutJson["proxy_ip"] = s5uSgwIp_up;
               OAMAGENT_LOG(INFO, "UserplanePatch s5uSgwIp_up (%s).\n", s5uSgwIp_up.c_str());				
           }
       }
       // get s1u configuration
       Json::Value s1uCfg = config.get ("s1u","Nil");
       if (0 == s1uCfg.compare("Nil")) {
           OAMAGENT_LOG(ERR, "[s1uCfg] is not found in request.\n");
    			
       }	
       else {
           string s1uIp_up  = s1uCfg.get ("up_ip_address","Nil").asString();			
           if (0 == s1uIp_up.compare("Nil")) {
               OAMAGENT_LOG(ERR, "[s1uIp_up] is not found in request.\n");
           }
	   else {
               sgwPutJson["s1u_ip"]   = s1uIp_up;
               // TEMP workaround
               sgwPutJson["s1u_nat_ip"]   = "";
               sgwPutJson["peer_ip"] = "192.190.111.111";
               //sgwPutJson["s1u_nat_ip"]   = s1uIp_up;				   
               OAMAGENT_LOG(INFO, "UserplanePatch s1u_ip (%s).\n", s1uIp_up.c_str());
            }
       }  
    }
    /////////////////start read selectors configuration ////////////////////////////////////////
    Json::Value selectors = request.get("selectors","Nil");
    if (0 == selectors.compare("Nil")) {
       OAMAGENT_LOG(ERR, "[selectors] is not found in request.\n");			
    }
    else {
       size = selectors.size();
       if (0 == size) {
           OAMAGENT_LOG(ERR, "[selectors] is not found.\n");
       }
       OAMAGENT_LOG(INFO, "UserplanePatch selectors size (%d).\n", size);
     
       for (int i = 0; i < size; i++) {
           /////// uli
           Json::Value uli = selectors[i].get ("uli","Nil");
           Json::Value tai = uli.get ("tai","Nil");
           #ifdef CUPS_API_INT64_TYPE
           int tac = tai.get ("tac","Nil").asInt();
           sgwPutJson["tac"] = utilsConvertDecIntToHexString(tac);
           OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] uli (tac %d).\n", i, tac);
           #else
           string tac = tai.get ("tac","Nil").asString();
           sgwPutJson["tac"] = tac;
           OAMAGENT_LOG(INFO, "UserplaneAdd selectors[%d] uli (tac %s).\n", i, tac.c_str());
           #endif
       }		
    }

    Json::Reader jsonReader;
    Json::FastWriter fastWriter;
    sgwPutData = fastWriter.write(sgwPutJson);
    OAMAGENT_LOG(INFO, "SGW_ADD JSON: %s\n",sgwPutJson.toStyledString().c_str());
    return 0;
}















