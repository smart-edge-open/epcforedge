/*******************************************************************************
This file contains the implementations of all the utility functions that I have
no idea where to put.
*******************************************************************************/
#include <string>
#include <fstream>
#include <streambuf>
//#include <iostream>

#include "TestUtility.h"

using std::string;
using std::ifstream;
using std::istreambuf_iterator;


string JSONFileToString(const string& file_name) {
    ifstream t{PATH_PREFIX + file_name + ".json"};
    string str {""};
    t.seekg(0, std::ios::end);	
    int size = t.tellg();
	if (-1 == size) {
        //return str;
        printf("JSONFileToString Failed with filename %s\n", file_name.c_str());
        throw -1;
    }
	
    str.reserve(size);
    t.seekg(0, std::ios::beg);			
    str.assign((istreambuf_iterator<char>(t)), istreambuf_iterator<char>());
    //std::cout << "json str := " << str << std::endl;
    return str;
}
