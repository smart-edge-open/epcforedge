/********************************************************************
 * SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2010-2014 Intel Corporation
 ********************************************************************/
/**
 * @file    Log.h
 * @brief   Header file for log declarations.
 ********************************************************************/

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
