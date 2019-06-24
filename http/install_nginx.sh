#! /bin/sh
############################################################################
# Copyright 2019 Intel Corporation. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
###########################################################################
setup_dir=${PWD}

set -e


log()
{
        green='\033[0;32m'
        reset='\e[0m'
        echo -e "${green}$1${reset}"
}


log "Install dependency package with yum"
yum -y install gcc
yum -y install openssl-devel.x86_64
yum -y install pcre-devel.x86_64


#Install Nginx
pkg_url=http://nginx.org/download/nginx-1.10.3.tar.gz
pkg_name=nginx-1.10.3
log "Download $pkg_name"
cd $setup_dir
wget -c $pkg_url
if [ $? -ne 0 ]; then
        log "$pkg_name package unavailable"
        exit 1
fi
log "Extract $pkg_name"
tar -xvmf $pkg_name.tar.gz > /dev/null
cd $setup_dir/$pkg_name
log "Build $pkg_name"
./configure --with-http_ssl_module && make install
if [ $? -ne 0 ]; then
        log "Compiled [ $pkg_name ] failed."
        exit 1
fi
# Link nginx executables to /usr/bin/
test -f /usr/bin/nginx && {
        /bin/cp -a /usr/bin/nginx /usr/bin/nginx-backup
        rm -rf /usr/bin/nginx
}
ln -s /usr/local/nginx/sbin/nginx /usr/bin/

# Link nginx configuration file to /etc/nginx/
test -d /etc/nginx || mkdir /etc/nginx
test -f /etc/nginx/nginx.conf && {
        /bin/cp -a /etc/nginx/nginx.conf /etc/nginx/nginx.conf-backup
        rm -rf /etc/nginx/nginx.conf
}
ln -s /usr/local/nginx/conf/nginx.conf /etc/nginx/

rm -rf $setup_dir/$pkg_name*

#Install libfcgi
pkg_url=ftp://ftp.linux.ro/gentoo/distfiles/fcgi-2.4.1-SNAP-0910052249.tar.gz
pkg_version=2.4.1-SNAP-0910052249
pkg_name=fcgi-2.4.1-SNAP-0910052249
log "Download $pkg_name"
cd $setup_dir
wget -c $pkg_url
if [ $? -ne 0 ]; then
        log "$pkg_name package unavailable"
        exit 1
fi
log "Extract $pkg_name"
tar -xvmf $pkg_name.tar.gz > /dev/null
cd $setup_dir/$pkg_name
log "Build $pkg_name"

sed -i 'N;24a#include <cstdio>' libfcgi/fcgio.cpp

./configure
make && make install
if [ $? -ne 0 ]; then
        log "Compiled [ $pkg_name ] failed."
        exit 1
fi
rm -rf $setup_dir/$pkg_name*


log "Done"



