//
// Created by david on 16-2-1.
//

#ifndef MECFCGI_TESTERFRAMEWORK_H
#define MECFCGI_TESTERFRAMEWORK_H

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


#endif //MECFCGI_TESTERFRAMEWORK_H
