/* SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2019 Intel Corporation
 */

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
