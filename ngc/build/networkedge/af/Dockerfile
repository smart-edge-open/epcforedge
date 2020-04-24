# SPDX-License-Identifier: Apache-2.0
# Copyright © 2019 Intel Corporation

FROM centos:7.6.1810

ENV http_proxy=$http_proxy
ENV https_proxy=$https_proxy
ENV no_proxy=localhost,127.0.0.1,afservice,nefservice,cntfservice

RUN yum install -y wget\
	net-tools \ 
	nmap \
        && yum clean all

WORKDIR /root
COPY ./af . 
COPY ./entrypoint.sh .

RUN chmod +x ./af
RUN chmod +x ./entrypoint.sh
CMD ["./entrypoint.sh"]

