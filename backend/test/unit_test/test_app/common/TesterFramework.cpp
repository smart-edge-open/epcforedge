//
// Created by david on 16-2-1.
//

#include "TesterFramework.h"

void TesterFramework::registerTest (TesterBase &tester, const string &desc)
{
  testers.push_back (pair<TesterBase *, string> (&tester, desc));
}

void TesterFramework::fireAllTests()
{
    for (pair<TesterBase *, string>tester : testers)
    {
        string additionalMessage;

        tester.first->execute (additionalMessage);
    }
}
