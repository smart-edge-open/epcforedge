#########################################################
# <COPYRIGHT_TAG>
#########################################################
cmake_minimum_required(VERSION 2.8)
project(oamagent)

#include_directories(
#        database
#        handler)

set(SOURCE_FILES
        main.cpp
        CpfInterface.cpp
        CupsMgmtInterface.cpp
        LocalConfig.cpp
        UpfController.cpp)

add_executable(oamagent ${SOURCE_FILES})
SET(EXECUTABLE_OUTPUT_PATH ${OAMAGENT_BINARY_DIR})

target_link_libraries(oamagent
        # oamagent libs
        liblog.a
        libexception.a
        libframework.a
        # system libs
        jsoncpp
        fcgi++
        fcgi
        curl 
        ssl
        crypto
        boost_system
        boost_thread
        pthread)
