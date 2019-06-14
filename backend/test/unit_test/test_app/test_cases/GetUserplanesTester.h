/********************************************************************
 * SPDX-License-Identifier: BSD-3-Clause
 * Copyright(c) 2010-2014 Intel Corporation
 ********************************************************************/
/*******************************************************************************
* Integration Tests for AppLiveIndicator, which is a handler for POST requests
* with a payload in JSON.
*******************************************************************************/

#ifndef UT_GETUSERPLANESTESTER_H
#define UT_GETUSERPLANESTESTER_H

#include <iostream>
#include <cstdlib>

#include "TesterBase.h"


class GetUserplanesTester: public TesterBase
{
    public:
    	int execute(string &additionalMessage);
};

#endif // #ifndef GET_H
