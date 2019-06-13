#########################################################
# <COPYRIGHT_TAG>
#########################################################
cmake_minimum_required(VERSION 2.8)
project(log)

aux_source_directory(. SOURCE_FILES)
add_library(log STATIC ${SOURCE_FILES})
SET(LIBRARY_OUTPUT_PATH ${OAMAGENT_LIBRARY_DIR})
