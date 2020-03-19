#!/bin/bash

#SPDX-License-Identifier: Apache-2.0
#Copyright © 2020 Intel Corporation

# Script to install curl 7.68.0. Sudo permissions are required to install it.
# Additional packages might be required, if script fails at some point install
# the following package and try it again.
#
# Additional packages might be required.
#   g++ make binutils autoconf automake autotools-dev libtool pkg-config 
#   zlib1g-dev libcunit1-dev libssl-dev libxml2-dev libev-dev libevent-dev 
#   libjansson-devi libjemalloc-dev cython python3-dev python-setuptools

set -x

#install nghttp2
sudo yum -y install nghttp2 libnghttp2 libnghttp2-devel

#install curl 7.68.0
wget https://curl.haxx.se/download/curl-7.68.0.tar.gz
tar -xf curl-7.68.0.tar.gz
rm curl-7.68.0.tar.gz
cd curl-7.68.0
./configure --with-nghttp2 --with-ssl
make -j

./src/curl --version
set +x
cd ../
echo -e "Curl 7.68.0 is installed in \033[0;32m `pwd`/src \033[0m"
echo -e "\033[0;33mrun ./curl-7.68.0/src/curl --version # it should print http2 as a feature\033[0m"
