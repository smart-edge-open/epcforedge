#########################################
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019 Intel Corporation
#########################################

FROM centos:7.6.1810

ENV http_proxy=$http_proxy
ENV https_proxy=$https_proxy
#RUN yum install -y libvirt-devel\
#        && yum clean all


WORKDIR /root/
COPY  nef_stub .
COPY ./entrypoint.sh .

CMD ["./entrypoint.sh"]

