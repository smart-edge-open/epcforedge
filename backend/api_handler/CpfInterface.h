/**
 * @file    cpfInterface.h
 * @brief   Header file of ControlPlane Function interface between  EPC OAMAgent and EPC control plane
 */

#ifndef __OAMAGENT__CPFINTERFACE__
#define __OAMAGENT__CPFINTERFACE__

#include <stdio.h>
#include <cstdint>
#include <iostream>
#include <memory>
#include <string>
#include <cstring>
#include <json/json.h>

using namespace std;

/* ------------------------------------------------------------------------- */
/* Public Function Defs */
/* ------------------------------------------------------------------------- */

int cpfCurlPost(string &url, string &postData, stringstream &responseData);
int cpfCurlGet(string &url, stringstream &responseData);
int cpfCurlDelete(string &url, bool &successFlg);
int cpfCurlPut(string &url, string &putData, stringstream &responseData);
int cpfCurlGetIdByItemIndex(int itemIndex, stringstream &responseData, string &upId);
int cpfCurlGetTacByItemIndex(int itemIndex, stringstream &responseData, string &tac);
int cpfCurlGetTotalCount(stringstream &responseData);
bool   cpfCurlGetSuccessFlag(stringstream &responseData);
string cpfCurlGetId(stringstream &responseData);
size_t cpfCurlDataHandle(void *ptr, size_t size, size_t nmemb, void *stream);


#endif

