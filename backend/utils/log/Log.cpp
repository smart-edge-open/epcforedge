/********************************************************************
 * SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2010-2014 Intel Corporation
 ********************************************************************/
/**
 * @file    Log.cpp
 * @brief   Implementation of Log.
 ********************************************************************/

#include "Log.h"
#include "iostream"
#include <stdint.h>
using namespace std;

void oamagentLogInit()
{
    const char *logid = "oamagent";
    openlog(logid, LOG_NDELAY | LOG_PID, LOG_DAEMON);
    return;
}
