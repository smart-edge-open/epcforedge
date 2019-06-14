#! /bin/sh
#########################################################
# <COPYRIGHT_TAG>
#########################################################
setup_dir=${PWD}

set -e


log()
{
        green='\033[0;32m'
        reset='\e[0m'
        echo -e "${green}$1${reset}"
}


log "Install dependency package with yum"
yum install -y cmake
yum -y install boost-devel.x86_64
yum -y install curl-devel

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


# Install Jsoncpp
pkg_url=https://github.com/open-source-parsers/jsoncpp.git
pkg_version=1.6.5
pkg_name=jsoncpp
log "Download $pkg_name"
cd $setup_dir
git clone -b $pkg_version $pkg_url
if [ $? -ne 0 ]; then
        log "$pkg_name package unavailable"
        exit 1
fi
cd $setup_dir/$pkg_name
log "Build $pkg_name"
cmake -H. -Bbuild && make -C build && make install -C build
if [ $? -ne 0 ]; then
        log "Compiled [ $pkg_name ] failed."
        exit 1
fi
rm -rf $setup_dir/$pkg_name*
cd $setup_dir

log "Done"



