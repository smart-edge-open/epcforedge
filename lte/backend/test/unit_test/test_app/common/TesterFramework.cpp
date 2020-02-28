/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */



#include "TesterFramework.h"

void TesterFramework::registerTest (TesterBase &tester, const string &desc)
{
  testers.push_back (pair<TesterBase *, string> (&tester, desc));
}

void TesterFramework::fireAllTests()
{
//    pair<TesterBase *, string> tester (NULL, "");
    for (const auto &tester : testers)
    {
        string additionalMessage;

        tester.first->execute (additionalMessage);
    }
}
