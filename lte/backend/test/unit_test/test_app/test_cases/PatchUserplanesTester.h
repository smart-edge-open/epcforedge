/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

/*******************************************************************************
* Integration Tests for PatchUserplane, which is a handler for POST requests
* with a payload in JSON.
*******************************************************************************/

#ifndef UT_PATCHUSERPLANES_H
#define UT_PATCHUSERPLANES_H

#include <iostream>
#include <cstdlib>

#include "TesterBase.h"


class PatchUserplanesTester: public TesterBase
{
    public:
    	int execute(string &additionalMessage);
};

#endif // #ifndef PATCH_H
