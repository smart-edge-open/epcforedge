/*******************************************************************************
* Integration Tests for AppLiveIndicator, which is a handler for POST requests
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

#endif // #ifndef MECFCGI_APPLIVEINDICATORTESTER_H
