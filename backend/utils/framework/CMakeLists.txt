#########################################################
# <COPYRIGHT_TAG>
#########################################################
cmake_minimum_required(VERSION 2.8)
project(framework)

aux_source_directory(. SOURCE_FILES)
add_library(framework STATIC ${SOURCE_FILES})
SET(LIBRARY_OUTPUT_PATH ${OAMAGENT_LIBRARY_DIR})

