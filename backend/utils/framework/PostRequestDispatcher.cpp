/************************************************************************************
 * <COPYRIGHT_TAG>
 ************************************************************************************/
/**
 * @file  PostRequestDispatcher.cpp
 * @brief POST method and JSON formatted request dispatcher.
 */

#include "PostRequestDispatcher.h"
#include "Exception.h"
#include "Log.h"

void PostRequestDispatcher::dispatchRequest(const string &action,
                                            Json::Value &request,
                                            Json::Value &response,
                                            map<string, string> &headers,
                                            map<string, string> &cookies)
{

	OAMAGENT_LOG(INFO, "PostRequestDispatcher with action %s.\n", action.c_str()); 	

    if (!action.length()) {
        OAMAGENT_LOG(ERR, "Dispatch failed.\n");
        throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
    }

    if (requestHandlers.find(action) != requestHandlers.end()) {
		OAMAGENT_LOG(INFO, "PostRequestDispatcher Find execute handler for the action %s.\n", action.c_str()); 	
        static_cast<PostRequestHandler *>(requestHandlers[action])->execute(request, response, headers, cookies);
        return;
    }
    OAMAGENT_LOG(ERR, "Dispatch failed, action: %s.\n", action.c_str());
    throw Exception(Exception::DISPATCH_NOTARGET, "Dispatch failed");
}

void PostRequestDispatcher::registerHandler(const string &action, PostRequestHandler &handler)
{
    requestHandlers[action] = &handler;
}
