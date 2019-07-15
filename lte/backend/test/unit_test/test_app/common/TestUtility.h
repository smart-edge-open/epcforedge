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
#ifndef TESTUTILITY_H
#define TESTUTILITY_H
#include <string>

using std::string;

// the path prefix for all json files, so that only the name of the file (minus
// its extension) needs to be passed to the function below.
const string PATH_PREFIX = "./json_payload/";


// reads a .json file in the path defined by PATH_PREFIX into a string
string JSONFileToString(const string& file_name);

#endif // TESTUTILITY_H
