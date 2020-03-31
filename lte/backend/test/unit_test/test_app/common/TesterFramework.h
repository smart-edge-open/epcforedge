/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */


#ifndef TESTERFRAMEWORK_H
#define TESTERFRAMEWORK_H

#include "TesterBase.h"
#include <deque>
#include <iostream>

using namespace std;
class TesterFramework
{
  deque<pair<TesterBase *, string>> testers;
public:
  void registerTest (TesterBase &tester, const string &desc);
  void fireAllTests();
};


#endif //TESTERFRAMEWORK_H
