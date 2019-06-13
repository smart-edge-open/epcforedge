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
 * @file    Log.h
 * @brief   Header file for log declarations.
 */

#ifndef __OAMAGENT__LOG__
#define __OAMAGENT__LOG__

#ifdef __cplusplus
extern "C" {
#endif

#include <stdio.h>
#include <syslog.h>
#include <string.h>

#define ERR      LOG_ERR
#define WARNING  LOG_WARNING
#define INFO     LOG_INFO

#define OUTPUT_MAX 20480
#define OAMAGENT_LOG(level,...) \
{ \
    char _fpt_log_buf[OUTPUT_MAX]; \
    snprintf(_fpt_log_buf, OUTPUT_MAX, __VA_ARGS__); \
    _fpt_log_buf[OUTPUT_MAX - 1] = 0; \
    syslog(level, "Func:%s(Line:%d)%s",  __FUNCTION__, __LINE__, _fpt_log_buf);\
}


//    printf("Func:%s(Line:%d)%s",  __FUNCTION__, __LINE__, _fpt_log_buf);\

void oamagentLogInit();

#ifdef __cplusplus
}
#endif /* extern "C" */
#endif /* defined(__OAMAGENT__LOG__) */
