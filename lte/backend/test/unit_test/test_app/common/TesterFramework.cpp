/*******************************************************************************
* Copyright 2019 Intel Corporation. All rights reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*******************************************************************************/


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
