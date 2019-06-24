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
 * @file    cpfInterface.h
 * @brief   Header file of ControlPlane Function interface between 
 *          EPC OAMAgent and EPC control plane
 ****************************************************************/

#ifndef __OAMAGENT__CPFINTERFACE__
#define __OAMAGENT__CPFINTERFACE__

#include <stdio.h>
#include <cstdint>
#include <iostream>
#include <memory>
#include <string>
#include <cstring>
#include <json/json.h>

using namespace std;

/* ------------------------------------------------------------------------- */
/* Public Function Defs */
/* ------------------------------------------------------------------------- */

int cpfCurlPost(string &url, string &postData, stringstream &responseData);
int cpfCurlGet(string &url, stringstream &responseData);
int cpfCurlDelete(string &url, bool &successFlg);
int cpfCurlPut(string &url, string &putData, stringstream &responseData);
int cpfCurlGetIdByItemIndex(int itemIndex, stringstream &responseData, string &upId);
int cpfCurlGetTacByItemIndex(int itemIndex, stringstream &responseData, string &tac);
int cpfCurlGetTotalCount(stringstream &responseData);
bool   cpfCurlGetSuccessFlag(stringstream &responseData);
string cpfCurlGetId(stringstream &responseData);
size_t cpfCurlDataHandle(void *ptr, size_t size, size_t nmemb, void *stream);


#endif

