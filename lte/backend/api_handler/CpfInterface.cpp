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
 * @file    cpfInterface.cpp
 * @brief   Implementation of ControlPlane Function interface 
 *          between EPC OAMAgent and EPC control plane
 **************************************************************/

#include <cstdint>
#include <iostream>
#include <memory>
#include <string>
#include <cstring>
#include <sstream>
#include <fcntl.h>
#include <sys/stat.h>

#include <curl/curl.h>
#include <json/json.h>
#include "CpfInterface.h"
#include "Log.h"

#define CPF_CURL_DEBUG 1

#if defined(UNIT_TEST) || defined(INT_TEST)
#include <string>
#include <fstream>
#include <streambuf>

int testUserplanesStart = 0;
const string PATH_PREFIX = "./json_payload/";
string JSONFileToString(const string& file_name) {
    ifstream t{PATH_PREFIX + file_name + ".json"};
    string str {""};
    t.seekg(0, std::ios::end);	
    int size = t.tellg();
	//printf("JSONFileToString  filename %s , size %d\n", file_name.c_str(), size);

	if (-1 == size) {
        //return str;
        printf("JSONFileToString Failed with filename %s\n", file_name.c_str());
        throw -1;
    }
	
    str.reserve(size);
    t.seekg(0, std::ios::beg);			
    str.assign((istreambuf_iterator<char>(t)), istreambuf_iterator<char>());
    //std::cout << "json str := " << str << std::endl;
    return str;
}
#endif

/**
* @brief			The utility function will be used to get field "success" value from reponse data of PEC
*					
* @param[in]		responseData 	The response data
* @return			The value of "success"
*/

bool cpfCurlGetSuccessFlag(stringstream &responseData)
{
    Json::Value  jsonData;
    Json::Reader jsonReader;
    jsonReader.parse(responseData.str().c_str(), jsonData);
    OAMAGENT_LOG(INFO, "Response: SuccessFlag %d\n",jsonData["success"].asBool());
    return jsonData["success"].asBool();	 
}

/**
* @brief			The utility function will be used to get field "id" value from reponse data of PEC
*					
* @param[in]		responseData 	The response data
* @return			The value of "id"
*/
string cpfCurlGetId(stringstream &responseData)
{
    Json::Value  jsonData;
    Json::Reader jsonReader;
    jsonReader.parse(responseData.str().c_str(), jsonData);
    // UserPlane id is just the index of the UPF configuration array
    // decimal value only
    // (now it is hex value, will fix later)
    OAMAGENT_LOG(INFO, "Response: Id %s\n",jsonData["id"].asString().c_str());
    return jsonData["id"].asString().c_str();	 
}

/**
* @brief			The utility function will be used to get field "id" value from reponse data of PEC
*					
* @param[in]		itemIndex 	    The array index.
* @param[in]		responseData 	The response data.
* @param[out]		upId 	        The "id" value.

* @return			Error Code
*/
int cpfCurlGetIdByItemIndex(int itemIndex, stringstream &responseData, string &upId)
{
    Json::Value  jsonData;
    Json::Reader jsonReader;
    jsonReader.parse(responseData.str().c_str(), jsonData);
    if(jsonData["items"][itemIndex]["id"].isString()){
        OAMAGENT_LOG(INFO, "Userplanes id %s\n",jsonData["items"][itemIndex]["id"].asString().c_str());		
        upId = jsonData["items"][itemIndex]["id"].asString();
        return 0;
    }
    else {
        OAMAGENT_LOG(ERR, "Userplanes id not exist for index %d.\n", itemIndex);	
        return -1;
    }
    return 0;	 
}

/**
* @brief			The utility function will be used to get field "TAC" value from reponse data of PEC
*					
* @param[in]		itemIndex 	    The array index.
* @param[in]		responseData 	The response data.
* @param[out]		tac 	        The "TAC" value.

* @return			Error Code
*/
int cpfCurlGetTacByItemIndex(int itemIndex, stringstream &responseData, string &tac)
{
    Json::Value  jsonData;
    Json::Reader jsonReader;
    jsonReader.parse(responseData.str().c_str(), jsonData);
    if(jsonData["items"][itemIndex]["tac"].isString()){
        OAMAGENT_LOG(INFO, "Userplanes tac %s\n",jsonData["items"][itemIndex]["tac"].asString().c_str());		
        tac = jsonData["items"][itemIndex]["tac"].asString();
        return 0;
    }
    else {
        OAMAGENT_LOG(ERR, "Userplanes tac not exist for index %d.\n", itemIndex);	
        return -1;
    }
    return 0;	 
}

/**
* @brief			The utility function will be used to get field "TotalCount" value from reponse data of PEC
*					
* @param[in]		responseData 	The response data
* @return			The value of "Total Count"
*/
int cpfCurlGetTotalCount(stringstream &responseData)
{
    Json::Value  jsonData;
    Json::Reader jsonReader;
    jsonReader.parse(responseData.str().c_str(), jsonData);
    if(!jsonData["totalCount"].isInt()){
		return -1;
    }

    OAMAGENT_LOG(INFO, "Response: totalCount %d\n",jsonData["totalCount"].asInt());
    return jsonData["totalCount"].asInt();	 
}

/**
* @brief			The callback function will be registered to CURL for the response data handling
*					
* @param[in]		ptr 	The pointer to the data
* @param[in]		size 	The data size
* @param[in]		nmemb   The num of memory block
* @param[out]		stream  The output data stream
* @return			Size of data
*/
size_t cpfCurlDataHandle(void *ptr, size_t size, size_t nmemb, void *stream) 
{
    string data((const char*) ptr, (size_t) size * nmemb);

    *((stringstream*) stream) << data << endl;

    return size * nmemb;
}

/**
* @brief			The inteface function will send POST to EPC and receive the response 
*					
* @param[in]		url 	      The URL for CURL post
* @param[in]		postData 	  The data to post
* @param[out]		responseData  The reponse
* @return			0:success; 1:failure
*/

int cpfCurlPost(string &url, string &postData, stringstream &responseData)
{
    CURL* curl;
    CURLcode ret;
    struct curl_slist* headers = NULL;
    string strResponseJson;
    long resCode;
    char strTemp[256] = "Content-Type: application/json";

    // Logging
    OAMAGENT_LOG(INFO, "Starting HTTP POST for url: %s\n", url.c_str());
#ifdef INT_TEST
    static int postTestCaseNum = 0;
    const char *postTestCaseRspData[2] = { 
         "{\"id\": \"5\",\"success\": true,\"msg\": \"\"}",  // PGW ADD
         "{\"id\": \"5\",\"success\": true,\"msg\": \"\"}"  // PGW ADD
    };

    //UT - direclty return
    responseData << postTestCaseRspData[postTestCaseNum];
    postTestCaseNum = (++postTestCaseNum)%2;
    return 0;
#endif

#ifdef UNIT_TEST
    static int postTestCaseNum = 0;
    if (postTestCaseNum == 7){
       // last test is negtive test
       // for post, mocked rsp data actually not enter
       return -1;
    }	
    const char *postTestCaseRspData[7] = { 
         "{\"id\": \"5\",\"success\": true,\"msg\": \"\"}",  // PGW ADD
	 "{\"id\": \"5\",\"success\": true,\"msg\": \"\"}",  // SGW ADD
	 "{\"id\": \"-1\",\"success\": false,\"msg\": \"\"}",  // PGW ADD Failed
	 "{\"id\": \"5\",\"success\": true,\"msg\": \"\"}",  // PGW ADD		
	 "{\"id\": \"-1\",\"success\": false,\"msg\": \"\"}",	 // SGW ADD Failed	
	 "{\"id\": \"5\",\"success\": true,\"msg\": \"\"}",  // PGW ADD with error ID response
	 "{\"id\": \"6\",\"success\": true,\"msg\": \"\"}"	 // SGW ADD with error ID response		
         };
		
    //UT - direclty return
    responseData << postTestCaseRspData[postTestCaseNum];
    postTestCaseNum++;
    return 0;
    #endif
	
    // curl init
    curl = curl_easy_init();
    if (curl == NULL) {
        OAMAGENT_LOG(ERR, "curl init failed\n");
        return -1;
    }

    // Set URL.
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_IPRESOLVE, CURL_IPRESOLVE_V4);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10);
    curl_easy_setopt(curl, CURLOPT_POST, 1); // set post flag
    
    // Set http headers
    headers = curl_slist_append(headers, strTemp);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    
    // Set post data
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, postData.c_str());
    curl_easy_setopt(curl, CURLOPT_POSTFIELDSIZE, postData.size());
        
    // Set reponse data hande
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, cpfCurlDataHandle);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &responseData);

    // Perform POST
    ret = curl_easy_perform(curl);
    if(ret != CURLE_OK) {
        OAMAGENT_LOG(ERR, "curl perform failed: %s\n",curl_easy_strerror(ret));
        return -1;        
    }
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &resCode);
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);


#if CPF_CURL_DEBUG
    // Print Response
    strResponseJson = responseData.str();
    printf("POST HTTP Response: code = %ld\n", resCode);
    printf("%s",strResponseJson.c_str());

    // Check http return code
    if (resCode == 200)
    {
        // DEUBG code
        printf("HTTP Response Success!\n");
    }
    else
    {
        OAMAGENT_LOG(ERR,"POST failed \n");
        return -1;
    }
#endif

    OAMAGENT_LOG(INFO, "Completed HTTP POST...\n");    
    return 0;
}

/**
* @brief			The inteface function will send GET to EPC and receive the response 
*					
* @param[in]		url 	      The URL for CURL get
* @param[out]		responseData  The reponse
* @return			0:success; 1:failure
*/
int cpfCurlGet(string &url, stringstream &responseData)
{
    CURL* curl;
    CURLcode ret;
    struct curl_slist* headers = NULL;
    string strResponseJson;
    long resCode;
    char strTemp[256] = "Content-Type: application/json";

    // Logging
    OAMAGENT_LOG(INFO, "Starting HTTP GET for url: %s\n", url.c_str());

#ifdef INT_TEST
    static int getTestCaseNum = 0;
    const char *getTestCaseRspData[4] = { 
		"PgwGetAllRspData",   // PGW GET ALL
		"SgwGetAllRspData",    // SGW GET ALL
		"PgwGetOneRspData",   // PGW GET ONE
		"SgwGetOneRspData"   // SGW GET ONE		
    };

    //UT - direclty return
    getTestCaseNum += testUserplanesStart; // get start index for the test cases
    responseData << JSONFileToString(getTestCaseRspData[getTestCaseNum]);	
    getTestCaseNum = (++getTestCaseNum)%2;
    //if (getTestCaseNum >= 2) getTestCaseNum = 0;	
    return 0;

#endif

#ifdef UNIT_TEST
    static int getTestCaseNum = 0;
    if (getTestCaseNum  == 20) {
       // last test is negtive test
       return -1;
    }	
    const char *getTestCaseRspData[20] = { 
		"PgwGetAllRspData",   // PGW GET ALL
		"SgwGetAllRspData",   // SGW GET ALL
		"PgwGetOneRspData",   // PGW GET ONE
		"SgwGetOneRspData",   // SGW GET ONE
		"PgwGetInvalidRspData",   // PGW GET Invalid
		"SgwGetInvalidRspData",   // SGW GET Invalid
		"PgwGetIdNotMatchRspData",   // PGW GET ID Invalid
		"SgwGetIdNotMatchRspData",   // SGW GET ID Invalid		
		"PgwGetTacNotMatchRspData",   // PGW GET TAC Invalid
		"SgwGetTacNotMatchRspData",   // SGW GET TAC Invalid
		"PgwGetTacNotFoundRspData",   // PGW GET TAC not found
		"SgwGetTacNotFoundRspData",   // SGW GET TAC not found	
		"PgwGetAPNWRONGRspData",   // PGW GET TAC not found
		"SgwGetAPNWRONGRspData",   // SGW GET TAC not found
		"PgwGetNoSgwS5uRspData",   // PGW GET TAC not found
		"SgwGetNoSgwS5uRspData",   // SGW GET TAC not found
		"PgwGetNoSgwTacRspData",   // PGW GET TAC not found
		"SgwGetNoSgwTacRspData",   // SGW GET TAC not found		
		"PgwGetNoSgwS1uRspData",   // PGW GET TAC not found
		"SgwGetNoSgwS1uRspData"    // SGW GET TAC not found				
		};

    //UT - direclty return
    responseData << JSONFileToString(getTestCaseRspData[getTestCaseNum]);	
    getTestCaseNum++;	
    return 0;
#endif

    
    // curl init
    curl = curl_easy_init();
    if (curl == NULL) {
        OAMAGENT_LOG(ERR, "curl init failed\n");
        return -1;
    }

    // Set URL.
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_IPRESOLVE, CURL_IPRESOLVE_V4);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10);
    
    // Set http headers
    headers=curl_slist_append(headers, strTemp);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        
    // Set reponse data hande
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, cpfCurlDataHandle);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &responseData);

    // Perform GET
    ret = curl_easy_perform(curl);
    if(ret != CURLE_OK) {
        OAMAGENT_LOG(ERR, "curl perform failed: %s\n",curl_easy_strerror(ret));
        return -1;        
    }
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &resCode);
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);


#if CPF_CURL_DEBUG
    // Print Response
    strResponseJson = responseData.str();
    printf("GET HTTP Response: code = %ld\n", resCode);
    printf("%s",strResponseJson.c_str());

    // Check http return code
    if (resCode == 200)
    {
        // DEUBG code
        printf("HTTP Response Success!\n");
    }
    else
    {
        OAMAGENT_LOG(ERR,"GET failed \n");
        return -1;
    }
#endif

    OAMAGENT_LOG(INFO, "Completed HTTP GET...\n");    
    return 0;
}

/**
* @brief			The inteface function will send DELETE to EPC and receive the response 
*					
* @param[in]		url 	      The URL for CURL delete
* @param[out]		successFlg    The operation success flag
* @return			0:success; 1:failure
*/
int cpfCurlDelete(string &url, bool &successFlg)
{
    CURL* curl;
    CURLcode ret;
    struct curl_slist* headers = NULL;
    string strResponseJson;
    long resCode;
    char strTemp[256] = "Content-Type: application/json";
    stringstream responseData;
    // Logging
    OAMAGENT_LOG(INFO, "Starting HTTP DELETE for url: %s\n", url.c_str());

    #ifdef INT_TEST
    static int delTestCaseNum = 0;
    const char *delTestCaseRspData[2] = { 
		"{\"success\":true,\"msg\":\"\"}",   //PGW respose for test 1
		"{\"success\":true,\"msg\":\"\"}"};   //SGW respose for test 1		
    static bool delTestCaseSucFlag[2] = {true, true};	
    //UT - direclty return
    successFlg   = delTestCaseSucFlag[delTestCaseNum];
    delTestCaseNum = (++delTestCaseNum)%2;	
    return 0;
    #endif

    #ifdef UNIT_TEST
    static int delTestCaseNum = 0;
    if (delTestCaseNum == 6) {
       // last test is negtive test
       return -1;
    }	
    const char *delTestCaseRspData[6] = { 
		"{\"success\":true,\"msg\":\"\"}",   //PGW respose for test 1
		"{\"success\":true,\"msg\":\"\"}",   //SGW respose for test 1			
		"{\"success\":true,\"msg\":\"\"}",   //PGW respose for test 1
		"{\"success\":true,\"msg\":\"\"}",   //SGW respose for test 1		
		"{\"success\":false,\"msg\":\"\"}",  //PGW respose for test 2
		"{\"success\":false,\"msg\":\"\"}"};  //SGW respose for test 2
    static bool delTestCaseSucFlag[6] = {true, true, true, true, false, false};	
    //UT - direclty return
    successFlg   = delTestCaseSucFlag[delTestCaseNum];
    delTestCaseNum++;	
    return 0;
    #endif // END UT TEST

	
    // curl init
    curl = curl_easy_init();
    if (curl == NULL) {
        OAMAGENT_LOG(ERR, "curl init failed\n");
        return -1;
    }

    // Set URL.
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_IPRESOLVE, CURL_IPRESOLVE_V4);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10);
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, "DELETE");
    
    // Set http headers
    headers=curl_slist_append(headers, strTemp);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        
    // Set reponse data hande
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, cpfCurlDataHandle);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &responseData);

    // Perform DELETE
    ret = curl_easy_perform(curl);
    if(ret != CURLE_OK) {
        OAMAGENT_LOG(ERR, "curl perform failed: %s\n",curl_easy_strerror(ret));
        return -1;        
    }
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &resCode);
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    successFlg = true; 
    printf("HTTP DELETE Response: code = %ld\n", resCode);
    OAMAGENT_LOG(INFO, "Completed HTTP DELETE...\n");    
    return 0;
}

/**
* @brief The inteface function will send PUT to EPC and receive the response 
*					
* @param[in]    url 	      The URL for CURL put
* @param[in]    putData       The data to put
* @param[out]   responseData  The reponse
* @return       0:success; 1:failure
*/
int cpfCurlPut(string &url, string &putData, stringstream &responseData)
{
    CURL* curl;
    CURLcode ret;
    struct curl_slist* headers = NULL;
    string strResponseJson;
    FILE * hd_src = NULL;
    struct stat file_info;
    long resCode;
    char strTemp[256] = "Content-Type: application/json";

    // Logging
    OAMAGENT_LOG(INFO, "Starting HTTP PUT for url: %s\n", url.c_str());
    // Generate json file for uploading
    hd_src = fopen("put_data.json", "wb");
    if (NULL == hd_src) {
       OAMAGENT_LOG(ERR, "open json file failed.\n");
       return -1; 
    }
    
    // write data into json file
    fwrite (putData.c_str() , sizeof(char), strlen(putData.c_str()), hd_src);
    fclose(hd_src);
    hd_src = NULL;
    
    /* get the file size of the local file */ 
    stat("put_data.json", &file_info);
    OAMAGENT_LOG(INFO, "Prepared json file size = %d\n",(int)file_info.st_size );   

#ifdef INT_TEST
    static int patchTestCaseNum = 0;
    const char *patchTestCaseRspData[2] = { 
		"{\"id\": \"5\",\"success\": true,\"msg\": \"\"}",  // PGW ADD
		"{\"id\": \"5\",\"success\": true,\"msg\": \"\"}"  // SGW ADD
		     };
		
    //UT - direclty return
    responseData << patchTestCaseRspData[patchTestCaseNum];
    patchTestCaseNum = (++patchTestCaseNum)%2;	
    return 0;
#endif

#ifdef UNIT_TEST
    static int patchTestCaseNum = 0;
    if (patchTestCaseNum == 3) {
        // last test is negtive test
	// for post, mocked rsp data actually not enter
	return -1;
    }

    const char *patchTestCaseRspData[4] = { 
		"{\"id\": \"5\",\"success\": true,\"msg\": \"\"}",  // PGW ADD
		"{\"id\": \"5\",\"success\": true,\"msg\": \"\"}",  // SGW ADD
		"{\"id\": \"-1\",\"success\": false,\"msg\": \"\"}",  // PGW ADD Failed
		"{\"id\": \"-1\",\"success\": false,\"msg\": \"\"}"	 // SGW ADD Failed	
		     };
		
    //UT - direclty return
    responseData << patchTestCaseRspData[patchTestCaseNum];
    patchTestCaseNum++;
    return 0;
#endif

    
    // curl init
    curl = curl_easy_init();
    if (curl == NULL) {
        OAMAGENT_LOG(ERR, "curl init failed\n");
        return -1;
    }

    // Set URL.
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_IPRESOLVE, CURL_IPRESOLVE_V4);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10);
    curl_easy_setopt(curl, CURLOPT_UPLOAD, 1L);
    curl_easy_setopt(curl, CURLOPT_PUT, 1L);
    hd_src = fopen("put_data.json", "rb");
    if (NULL == hd_src) {
       OAMAGENT_LOG(ERR, "open json file failed.\n"); 
       return -1;  
    }    
    curl_easy_setopt(curl, CURLOPT_READDATA, hd_src);
    curl_easy_setopt(curl, CURLOPT_INFILESIZE_LARGE,
                     (curl_off_t)file_info.st_size);    
    
    // Set http headers
    headers=curl_slist_append(headers, strTemp);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    
        
    // Set reponse data hande
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, cpfCurlDataHandle);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &responseData);

    // Perform PUT
    ret = curl_easy_perform(curl);
    if(ret != CURLE_OK) {
        OAMAGENT_LOG(ERR, "curl perform failed: %s\n",curl_easy_strerror(ret));
        fclose(hd_src);
        return -1;        
    }
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &resCode);
    curl_slist_free_all(headers);
    curl_easy_cleanup(curl);
    fclose(hd_src);

#if CPF_CURL_DEBUG
    // Print Response
    strResponseJson = responseData.str();
    printf("PUT HTTP Response: code = %ld\n", resCode);
    printf("%s",strResponseJson.c_str());

    // Check http return code
    if (resCode == 200)
    {
        // DEUBG code
        printf("HTTP Response Success!\n");
    }
    else
    {
        OAMAGENT_LOG(ERR,"PUT failed \n");
        return -1;
    }
#endif

    OAMAGENT_LOG(INFO, "Completed HTTP PUT...\n");    
    return 0;
}
























