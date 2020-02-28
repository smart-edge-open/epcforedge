/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

#ifndef UT_DELUSERPLANESTESTER_H
#define UT_DELUSERPLANESTESTER_H

#include <iostream>
#include <cstdlib>

#include "TesterBase.h"


class DelUserplanesTester: public TesterBase
{
    public:
    	int execute(string &additionalMessage);
};

#endif // #ifndef DEL
