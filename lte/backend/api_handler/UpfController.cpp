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
 * @file  UpfController.cpp
 * @brief OAMAGENT handler UpfControl Message from MEC controller
 ********************************************************************/

#include "UpfController.h"
#include "Exception.h"
#include "HandlerCommon.h"
#include "Log.h"
#include "RawRequest.h"
#include "LocalConfig.h"
#include "CpfInterface.h"
#include "CupsMgmtInterface.h"


#ifdef INT_TEST
extern int testUserplanesStart;
#endif
/**
* @brief                Post new user plane configuration
* @param[in]		request 	JSON-formatted request data.
* @param[out]		response	JSON-formatted key-value pair(s) indicating
*								response.
* @param[out]		headers 	Response headers.
* @param[in]		cookies 	Cookies header in request.
* @throw		Exception	Thrown on failure.
* @return		void
*/

void UserplaneAdd::execute(Json::Value &request, Json::Value &response,
                    map<string, string> &headers, map<string, string> &cookies)
{
    try {
        string pgwGetUrl, sgwGetUrl;
        string pgwPostData,sgwPostData;
	stringstream pgwPostResponse,sgwPostResponse;
        string pgwId, sgwId;
		
	// Prepare PGW and SGW URL
        pgwGetUrl = "http://" + localcfg_pgw_ipaddress + ":" + localcfg_pgw_port + \
				 "/api/v1/pgwprofile?entity-type=pgw-dpf";
        sgwGetUrl = "http://" + localcfg_sgw_ipaddress + ":" + localcfg_sgw_port + \
				 "/api/v1/sgwprofile?entity-type=sgw-dpf";
			
        // Check function exist
	string function = request.get("function", "Nil").asString();
	if (0 == function.compare("Nil") || 
	  ((0 != function.compare("SAEGWU")) && (0 != function.compare("PGWU")) && (0 != function.compare("SGWU")))) {
            OAMAGENT_LOG(ERR, "[function] is not found in request.\n");
            throw Exception(Exception::INVALID_UERPLANE_FUNCTION);
	}
        OAMAGENT_LOG(INFO, "UserplaneAdd execute with function (%s).\n", function.c_str());	        

        // POST PGW
        if ((0 == function.compare("SAEGWU")) || (0 == function.compare("PGWU"))) {			

            // Prepare postdata
            if (0 != cupsMgmtMsg.fillPostPgwRequest(request, pgwPostData)) {
				OAMAGENT_LOG(ERR, "filling message falied.\n");
				throw Exception(Exception::INVALID_UERPLANE_FUNCTION);
            }
			
            // Post PGW config 
            if (0 != cpfCurlPost(pgwGetUrl,pgwPostData,pgwPostResponse)) {
                OAMAGENT_LOG(ERR, "curl POST failed.\n");				
                throw Exception(Exception::CONNECT_EPC_ERROR);
            }

            // Check Operation Success Flag
            if (false == cpfCurlGetSuccessFlag(pgwPostResponse)) {
                OAMAGENT_LOG(ERR," PgwPostResponse Success Flag is False \n");
                throw Exception(Exception::ADDED_USERPLANE);
            }

            // Get ID from repsonse
            pgwId = cpfCurlGetId(pgwPostResponse);
            response["id"] = pgwId.c_str();			
			
        }	
		
        // SGWU only
        if ((0 == function.compare("SAEGWU")) || (0 == function.compare("SGWU"))) {

            // Prepare postdata
            if (0 != cupsMgmtMsg.fillPostSgwRequest(request, sgwPostData)) {
                 OAMAGENT_LOG(ERR, "filling message falied.\n");
                 throw Exception(Exception::INVALID_UERPLANE_FUNCTION);
            }

            // Post SGW config 
            if (0 != cpfCurlPost(sgwGetUrl,sgwPostData,sgwPostResponse)) {
                OAMAGENT_LOG(ERR, "curl POST failed.\n");				
                throw Exception(Exception::CONNECT_EPC_ERROR);
            }
            // Check Operation Success Flag
            if (false == cpfCurlGetSuccessFlag(sgwPostResponse)) {
                OAMAGENT_LOG(ERR," PgwPostResponse Success Flag is False \n");
                throw Exception(Exception::ADDED_USERPLANE);
            }

            // Get ID from repsonse
            sgwId = cpfCurlGetId(sgwPostResponse);
            response["id"] = sgwId.c_str();			
			
        }

        // check SGW ID and PGW ID
        if (0 == function.compare("SAEGWU")) {	
            if (0 != strcmp(sgwId.c_str(), pgwId.c_str())) {
                OAMAGENT_LOG(ERR," SGW and PGW ID not match \n");
                throw Exception(Exception::INTERNAL_SOFTWARE_ERROR);
            }
        }

        headers["Status"] = HTTP_SC_OK;
        response["result"] = "OK";

    }
    catch (Exception &e) {
        string res;
        string statusCode;
        Exception::handlerException(e, res, statusCode);
        headers["Status"] = statusCode;
        response["result"] = res;
    }
}


/**
* @brief		Update the configuration of specific user plane
* @param[in]		request 	JSON-formatted request data.
* @param[out]		response	JSON-formatted key-value pair(s) indicating
*					response.
* @param[out]		headers 	Response headers.
* @param[in]		cookies 	Cookies header in request.
* @throw		Exception	Thrown on failure.
* @return		void
*/

void UserplanePatchByID::execute(Json::Value &request, Json::Value &response,
                    map<string, string> &headers, map<string, string> &cookies)
{
   try {
        
        OAMAGENT_LOG(INFO, "UserplanePatchByID(%s) Executing.\n", request["UUID"].asString().c_str());

        string pgwGetUrl, sgwGetUrl;
        string pgwPostData,sgwPostData;
        string pgwId,sgwId;
        stringstream pgwPostResponse,sgwPostResponse;
 
        // Prepare PGW and SGW URL
        pgwGetUrl = "http://" + localcfg_pgw_ipaddress + ":" + localcfg_pgw_port + \
			 "/api/v1/pgwprofile?entity-type=pgw-dpf&id=" + request["UUID"].asString();
        sgwGetUrl = "http://" + localcfg_sgw_ipaddress + ":" + localcfg_sgw_port + \
			 "/api/v1/sgwprofile?entity-type=sgw-dpf&id=" + request["UUID"].asString();
		 				
        // Check function exist
        string function = request.get("function", "Nil").asString();
        if (0 == function.compare("Nil")) {
            OAMAGENT_LOG(ERR, "[function] is not found in request.\n");
            throw Exception(Exception::INVALID_UERPLANE_FUNCTION);
        }

        // Check function value		
        OAMAGENT_LOG(INFO, "UserplaneAdd execute with function (%s).\n", function.c_str()); 	
        if (0 == function.compare("NONE")) {
            OAMAGENT_LOG(ERR, "[function] is not found in request.\n");
            throw Exception(Exception::INVALID_UERPLANE_FUNCTION);
        }

        // Combined PGWU and SGW
        if ((0 == function.compare("SAEGWU")) || (0 == function.compare("PGWU"))) {	
           // Prepare postdata
           if (0 != cupsMgmtMsg.fillPutPgwRequest(request, pgwPostData)) {
              OAMAGENT_LOG(ERR, "filling message falied.\n");
              throw Exception(Exception::INTERNAL_SOFTWARE_ERROR);
           }

           // PUT PGW config 
           if (0 != cpfCurlPut(pgwGetUrl,pgwPostData,pgwPostResponse)) {
              OAMAGENT_LOG(ERR, "curl PUT failed.\n");				
              throw Exception(Exception::CONNECT_EPC_ERROR);
           }

           // Check Operation Success Flag
           if (false == cpfCurlGetSuccessFlag(pgwPostResponse)) {
              OAMAGENT_LOG(ERR," PgwPatchResponse Success Flag is False \n");
              throw Exception(Exception::USERPLANE_NOT_FOUND);
           }
			
        }		

        // SGWU only
        if ((0 == function.compare("SAEGWU")) || (0 == function.compare("SGWU"))) {		
           // Prepare postdata
           if (0 != cupsMgmtMsg.fillPutSgwRequest(request, sgwPostData)) {
              OAMAGENT_LOG(ERR, "filling message falied.\n");
              throw Exception(Exception::INTERNAL_SOFTWARE_ERROR);
           }

           // Post SGW config 
           if (0 != cpfCurlPut(sgwGetUrl,sgwPostData,sgwPostResponse)) {							
              OAMAGENT_LOG(ERR, "curl PUT failed.\n");				
              throw Exception(Exception::CONNECT_EPC_ERROR);
           }

           // Check Operation Success Flag
           if (false == cpfCurlGetSuccessFlag(sgwPostResponse)) {
              OAMAGENT_LOG(ERR," SgwPatchResponse Success Flag is False \n");
              throw Exception(Exception::USERPLANE_NOT_FOUND);
           }
        }
		
        headers["Status"] = HTTP_SC_OK;
        response["result"] = "OK";

    }
    catch (Exception &e) {
        string res;
        string statusCode;
        Exception::handlerException(e, res, statusCode);
        headers["Status"] = statusCode;
        response["result"] = res;
    }
}


/**
 * @brief            Get User Planes List.
 * @param[in]        request     JSON-formatted request data.
 * @param[out]       response    JSON-formatted key-value pair(s) indicating
 *                               response.
 * @param[out]       headers     Response headers.
 * @param[in]        cookies     Cookies header in request.
 * @throw            Exception   Thrown on failure.
 * @return           void
 */

void UserplanesListGet::execute(map<string, string> params, Json::Value &response,
                    map<string, string> &headers, map<string, string> &cookies)
{
    try {
        OAMAGENT_LOG(INFO, "UserplanesListGet Executing.\n");

        // Preparing URL for get 
        int pgwCount, sgwCount, itemIndex;
        stringstream pgwGetData, sgwGetData;        
        string pgwGetUrl = "http://" + localcfg_pgw_ipaddress + ":" + localcfg_pgw_port + \
                     "/api/v1/pgwprofile?action=list&entity-type=pgw-dpf";		
        string sgwGetUrl = "http://" + localcfg_sgw_ipaddress + ":" + localcfg_sgw_port + \
                     "/api/v1/sgwprofile?action=list&entity-type=sgw-dpf";       

        #ifdef INT_TEST
        testUserplanesStart = 0; // start index for the user plane mock json data
        #endif

        // Get PGW information from CP
        if (0 != cpfCurlGet(pgwGetUrl, pgwGetData)) {
            throw Exception(Exception::CONNECT_EPC_ERROR);
        }
        // Get SGW information from CP
        if (0 != cpfCurlGet(sgwGetUrl, sgwGetData)) {
            throw Exception(Exception::CONNECT_EPC_ERROR);
        }

        Json::Value jsonPgwData;     		
        Json::Value jsonSgwData;
        Json::Reader jsonReader;
        Json::Value  responseItem[MAX_USERPLANES_NUM];
		
        // Check Total Account
        if (0 > (pgwCount = cpfCurlGetTotalCount(pgwGetData))) {
            OAMAGENT_LOG(ERR, "Wrong PGW totalCount %d.\n", pgwCount);
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }
		
        // Check Total Account
        if (0 > (sgwCount = cpfCurlGetTotalCount(sgwGetData))) {
            OAMAGENT_LOG(ERR, "Wrong SGW totalCount %d.\n", sgwCount);
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }	   

        if (pgwCount != sgwCount) {
            OAMAGENT_LOG(ERR, "Wrong SGW PGW totalCount %d - %d.\n", sgwCount, pgwCount);
            throw Exception(Exception::INTERNAL_SOFTWARE_ERROR);			
        }
        OAMAGENT_LOG(INFO, "Preparing response.\n");
        jsonReader.parse(pgwGetData.str().c_str(), jsonPgwData);
        jsonReader.parse(sgwGetData.str().c_str(), jsonSgwData);

        for (itemIndex = 0; itemIndex < pgwCount; itemIndex++) {
          // Fill Response
           if (0 != cupsMgmtMsg.fillGetPgwResponse(jsonPgwData, itemIndex, responseItem[itemIndex])) {
              throw Exception(Exception::INTERNAL_SOFTWARE_ERROR);		   	   
           }
           if (0 != cupsMgmtMsg.fillGetSgwResponse(jsonSgwData, itemIndex, responseItem[itemIndex])) {
              throw Exception(Exception::INTERNAL_SOFTWARE_ERROR);
           }
           response["userplanes"].append(responseItem[itemIndex]);
        }
	
        OAMAGENT_LOG(INFO, "UserplanesListGet Success With count (%d).\n", pgwCount);
        headers["Status"] = HTTP_SC_OK;
        response["result"] = "OK";

    }
    catch (Exception &e) {
        string res;
        string statusCode;
        Exception::handlerException(e, res, statusCode);
		OAMAGENT_LOG(ERR, "UserplanesListGet Failed (%s).\n", statusCode.c_str());				
        headers["Status"] = statusCode;
        response["result"] = res;
    }
}

/**
 * @brief			Get Specific User Plane.
 * @param[in]		request 	JSON-formatted request data.
 * @param[out]		response	JSON-formatted key-value pair(s) indicating
 *								response.
 * @param[out]		headers 	Response headers.
 * @param[in]		cookies 	Cookies header in request.
 * @throw			Exception	Thrown on failure.
 * @return			void
 */
void UserplaneGetByID::execute(map<string, string> params, Json::Value &response,
                    map<string, string> &headers, map<string, string> &cookies)
{
    try {
        OAMAGENT_LOG(INFO, "UserplaneGetByID(%s) Executing.\n", params["UUID"].c_str());

        // Preparing URL for get 
        int pgwCount, sgwCount, itemIndex;
        string pgwId, sgwId;
        string pgwTac, sgwTac;
        stringstream pgwGetData,sgwGetData;        
        string pgwGetUrl = "http://" + localcfg_pgw_ipaddress + ":" + localcfg_pgw_port + \
                    "/api/v1/pgwprofile?action=list&entity-type=pgw-dpf&id=" + params["UUID"];
        string sgwGetUrl = "http://" + localcfg_sgw_ipaddress + ":" + localcfg_sgw_port + \
                    "/api/v1/sgwprofile?action=list&entity-type=sgw-dpf&id=" + params["UUID"]; 

        #ifdef INT_TEST
        testUserplanesStart = 2; // start index for the user plane mock json data
        #endif



        Json::Value jsonPgwData;     		
        Json::Value jsonSgwData;
        Json::Reader jsonReader;
        itemIndex = 0;
		
        // Get PGW information from CP
        if (0 != cpfCurlGet(pgwGetUrl, pgwGetData)) {
            OAMAGENT_LOG(ERR, "Get PGW Failed \n");
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }
		
        // Get SGW information from CP
        if (0 != cpfCurlGet(sgwGetUrl, sgwGetData)) {
            OAMAGENT_LOG(ERR, "Get SGW Failed \n");
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }

        // Check Total Account
        if (1 != (pgwCount = cpfCurlGetTotalCount(pgwGetData))) {
            OAMAGENT_LOG(ERR, "Wrong PGW totalCount %d.\n", pgwCount);
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }
		
        // Check Total Account
        if (1 != (sgwCount = cpfCurlGetTotalCount(sgwGetData))) {
            OAMAGENT_LOG(ERR, "Wrong SGW totalCount %d.\n", sgwCount);
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }

        // Get PGW ID and SGW ID, then check whether it match
        if (0 != cpfCurlGetIdByItemIndex(itemIndex,pgwGetData,pgwId)
          ||0 != cpfCurlGetIdByItemIndex(itemIndex,sgwGetData,sgwId)) {
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }

		
        OAMAGENT_LOG(INFO, "PGW id(%s) or SGW Id (%s) for GetID (%s).\n",
				pgwId.c_str(),sgwId.c_str(),params["UUID"].c_str());		
        #if 1 // need check when value is "09", return ID = "9". 
        if (0 != strcmp(sgwId.c_str(), params["UUID"].c_str())
		 || 0 != strcmp(pgwId.c_str(), params["UUID"].c_str())) { // check result
            OAMAGENT_LOG(ERR, "Not valid PGW id(%s) or SGW Id (%s) for GetID (%s).\n",
				pgwId.c_str(),sgwId.c_str(),params["UUID"].c_str());
            throw Exception(Exception::USERPLANE_NOT_FOUND);;
    	}
        #endif

        // Get TAC, then check TAC validation
        if (0 != cpfCurlGetTacByItemIndex(itemIndex,pgwGetData,pgwTac)
         || 0 != cpfCurlGetTacByItemIndex(itemIndex,sgwGetData,sgwTac)) {
            throw Exception(Exception::INTERNAL_SOFTWARE_ERROR);
        }
		
        if (0 != strcmp(sgwTac.c_str(), pgwTac.c_str())) { // check result
            OAMAGENT_LOG(ERR, "Not valid PGW TAC (%s) or SGW TAC (%s) for GetID (%s).\n",
				pgwTac.c_str(),sgwTac.c_str(),params["UUID"].c_str());
            throw Exception(Exception::INTERNAL_SOFTWARE_ERROR);;
    	}		

        // Fill Response
        OAMAGENT_LOG(INFO, "Preparing response.\n");
        jsonReader.parse(pgwGetData.str().c_str(), jsonPgwData);
        jsonReader.parse(sgwGetData.str().c_str(), jsonSgwData);
        if (0!= cupsMgmtMsg.fillGetPgwResponse(jsonPgwData, itemIndex, response)){
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }
        if (0 != cupsMgmtMsg.fillGetSgwResponse(jsonSgwData, itemIndex, response)){
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }
		OAMAGENT_LOG(INFO, "UserplanesGet For (%s) Success.\n",params["UUID"].c_str());
        headers["Status"] = HTTP_SC_OK;
        response["result"] = "OK";

    }
    catch (Exception &e) {
		
        string res;
        string statusCode;
        Exception::handlerException(e, res, statusCode);
		OAMAGENT_LOG(ERR, "UserplanesGet For (%s) Failed (%s).\n", params["UUID"].c_str(),statusCode.c_str());		
        headers["Status"] = statusCode;
        response["result"] = res;
    }
}

/**
 * @brief            Delete specific user plane.
 * @param[in]        params      JSON-formatted params data.
 * @param[out]       response    JSON-formatted key-value pair(s) indicating
 *                               response.
 * @param[out]       headers     Response headers.
 * @param[in]        cookies     Cookies header in request.
 * @throw            Exception   Thrown on failure.
 * @return           void
 */

void UserplaneDelByID::execute(map<string, string> params,
                                Json::Value &response,
                                map<string, string> &headers,
                                map<string, string> &cookies)
{
    try {

        OAMAGENT_LOG(INFO, "UserplaneDelByID(%s) Executing.\n",params["UUID"].c_str());
        string pgwGetUrl = "http://" + localcfg_pgw_ipaddress + ":" + localcfg_pgw_port + \
                     "/api/v1/pgwprofile?entity-type=pgw-dpf&id=" + params["UUID"];
        string sgwGetUrl = "http://" + localcfg_sgw_ipaddress + ":" + localcfg_sgw_port + \
                     "/api/v1/sgwprofile?entity-type=sgw-dpf&id=" + params["UUID"];

        // Delete PGW information from CP
        bool sucFlag = false;
        if (0 != cpfCurlDelete(pgwGetUrl, sucFlag)) {
            OAMAGENT_LOG(ERR, "DeletePGW failed.\n");
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }
		
        if (sucFlag == false) {
	        OAMAGENT_LOG(ERR, "DeletePGW failed.\n");
            throw Exception(Exception::USERPLANE_NOT_FOUND);			
        }
		
        // Delete SGW information from CP
        if (0 != cpfCurlDelete(sgwGetUrl, sucFlag)) {
            OAMAGENT_LOG(ERR, "DeleteSGW failed.\n");
            throw Exception(Exception::USERPLANE_NOT_FOUND);
        }
		
        if (sucFlag == false) {
	        OAMAGENT_LOG(ERR, "DeletePGW failed.\n");
            throw Exception(Exception::USERPLANE_NOT_FOUND);			
        }		

        headers["Status"] = HTTP_SC_OK;
        response["result"] = "OK";
    }
    catch (Exception &e) {
        string res;
        string statusCode;
        Exception::handlerException(e, res, statusCode);
        headers["Status"] = statusCode;
        response["result"] = res;
    }
}


