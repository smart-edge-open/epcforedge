/************************************************************************************
 * <COPYRIGHT_TAG>
 ************************************************************************************/
/**
 * @file  GetRequestDispatcher.cpp
 * @brief GET method requests dispatcher
 */

#include "GetRequestDispatcher.h"
#include "Exception.h"
#include "Log.h"

#define REQUEST_SPLIT_MARK         "/"
#define REQUEST_SPLIT_MARK_LENGTH  1
void GetRequestDispatcher::dispatchRequest(const string &action,
                                            map<string, string> &params,
                                            Json::Value &response,
                                            map<string, string> &headers,
                                            map<string, string> &cookies)
{
	OAMAGENT_LOG(INFO, "GetRequestDispatcher with xxx action %s.\n", action.c_str()); 	

    if (!action.length()) {
        OAMAGENT_LOG(ERR, "Dispatch failed.\n");
        throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
    }
    //string action_bk = action;
    size_t pos;
    if (requestHandlers.find(action) != requestHandlers.end()) {
    	OAMAGENT_LOG(INFO, "GetRequestDispatcher Find execute handler for the action (%s).\n", action.c_str()); 			
        static_cast<GetRequestHandler *>(requestHandlers[action])->execute(params, response, headers, cookies);
        return;
    } else if ((pos = action.find_last_of(REQUEST_SPLIT_MARK)) != string::npos) {
    	OAMAGENT_LOG(INFO, "GetRequestDispatcher Find SplitMask for the action %s.\n", action.c_str());
        string newAction = action.substr(0, pos) + "/UUID";
        if (requestHandlers.find(newAction) != requestHandlers.end()) {
    	    //OAMAGENT_LOG(INFO, "GetRequestDispatcher substr(%d,%d) with pos =%d, actlen=%d\n",
			//	pos + REQUEST_SPLIT_MARK_LENGTH, action.length() - pos + REQUEST_SPLIT_MARK_LENGTH, pos, action.length());			
            params["UUID"] = action.substr(pos + REQUEST_SPLIT_MARK_LENGTH, action.length() - pos + REQUEST_SPLIT_MARK_LENGTH);
    	    OAMAGENT_LOG(INFO, "GetRequestDispatcher Find UUID (%s) for the newaction (%s)\n",params["UUID"].c_str(), newAction.c_str());
			if (0 == strlen(params["UUID"].c_str())) {
				throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
			}
			//if (0 == strlen(params["UUID"].c_str())) {
			  // no UUID for the URL. so use old action
		    //OAMAGENT_LOG(INFO, "len: (%d) \n", params["UUID"].length());	  
			  //static_cast<GetRequestHandler *>(requestHandlers[newAction])->execute(params, response, headers, cookies);
			//}
			//else {
              static_cast<GetRequestHandler *>(requestHandlers[newAction])->execute(params, response, headers, cookies);
			//}
            return;
			
        }
    }
    OAMAGENT_LOG(ERR, "Dispatch failed, action: %s.\n", action.c_str());
    throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
}

void GetRequestDispatcher::registerHandler(const string &action, GetRequestHandler &handler)
{
    requestHandlers[action] = &handler;
}
