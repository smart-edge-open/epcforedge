/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/*******************************************************************************
* Integration Tests for PostUserplane, which is a handler for POST requests
* with a payload in JSON.
*******************************************************************************/

#ifndef UT_POSTUSERPLANES_H
#define UT_POSTUSERPLANES_H

#include <iostream>
#include <cstdlib>

#include "TesterBase.h"


class PostUserplanesTester: public TesterBase
{
    public:
    	int execute(string &additionalMessage);
};

#endif // #ifndef POST_H
