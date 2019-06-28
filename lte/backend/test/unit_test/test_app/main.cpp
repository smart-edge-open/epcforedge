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
//
//

#include "TesterFramework.h"

#include "PostUserplanesTester.h"
#include "PatchUserplanesTester.h"

#include "DelUserplanesTester.h"
#include "GetUserplanesTester.h"
extern int okNum ;
extern int ngNum ;

int main()
{
    // declarations here
    TesterFramework framework;

    // test cases


    printf("[==========] Starting register tester\n");
    PostUserplanesTester post_userplanes_tester;
    framework.registerTest(post_userplanes_tester, "PostUserplanes Test");

    PatchUserplanesTester patch_userplanes_tester;
    framework.registerTest(patch_userplanes_tester, "PatchUserplanes Test");
	

    GetUserplanesTester get_userplanes_tester;
    framework.registerTest(get_userplanes_tester, "GetUserplanes Test");


    DelUserplanesTester del_userplanes_tester;
    framework.registerTest(del_userplanes_tester, "DelUserplanes Test");

             
    printf("[==========] Running tester\n");

    // test execution starts here
    framework.fireAllTests();
    printf("\n[==========] Completed tester: OKNum=%d, TotalNum=%d\n", okNum, okNum + ngNum);

    return 0;
}
