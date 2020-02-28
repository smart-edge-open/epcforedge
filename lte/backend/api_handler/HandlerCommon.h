/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/**
 * @file    HandlerCommon.h
 * @brief   Header file for common declarations.
 ********************************************************************/

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

const uint32_t IP_MASK_MIN;
const uint32_t IP_MASK_MAX = 32;
const uint32_t PORT_MIN;
const uint32_t PORT_MAX = 65535;

const uint32_t SERVICE_ID_LENGTH = 64;


#endif /* defined(__OAMAGENT__HANDLERCOMMON__) */
