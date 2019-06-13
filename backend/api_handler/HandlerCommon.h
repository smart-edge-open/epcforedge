/************************************************************************************
Copyright 2019 Intel Corporation. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

************************************************************************************/
/**
 * @file    HandlerCommon.h
 * @brief   Header file for common declarations.
 */

#ifndef __OAMAGENT__HANDLERCOMMON__
#define __OAMAGENT__HANDLERCOMMON__

#include <string>
#include <stdlib.h>
#include <json/json.h>

using namespace std;

const char SPLIT_MARK = ':';
const char SPLIT_MARKL = '<';
const char SPLIT_MARKR = '>';

const string SCHEMA_KEYFIELDS = "keyFields";
const string SCHEMA_DATA = "data";



const char RULE_FIELD_SPLIT_MARK = ',';
const char RULE_FIELD_VALUE_SPLIT_MARK = ':';
const char RULE_FIELD_VALUE_EQUAL_MARK = '=';
const char RULE_FIELD_MIN_MIN_SPLIT_MARK = '-';
const char RULE_FIELD_IP_MASK_SPLIT_MARK = '/';

const uint32_t IP_MASK_MIN = 0;
const uint32_t IP_MASK_MAX = 32;
const uint32_t PORT_MIN = 0;
const uint32_t PORT_MAX = 65535;

const uint32_t SERVICE_ID_LENGTH = 64;


#endif /* defined(__OAMAGENT__HANDLERCOMMON__) */
