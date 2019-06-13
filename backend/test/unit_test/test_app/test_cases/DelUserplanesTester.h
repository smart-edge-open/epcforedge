/*******************************************************************************
* Integration Tests for AppLiveIndicator, which is a handler for POST requests
* with a payload in JSON.
*******************************************************************************/

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

#endif // #ifndef MECFCGI_APPLIVEINDICATORTESTER_H
