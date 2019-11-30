# Copyright 2019 Intel Corporation, Inc. All rights reserved
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

FROM centos:7.6.1810

ENV http_proxy=$http_proxy
ENV https_proxy=$https_proxy
ENV no_proxy=localhost,127.0.0.1

RUN yum install -y wget\
	net-tools \ 
	nmap \
        && yum clean all

WORKDIR /root
COPY ./nef . 
COPY ./entrypoint.sh .

RUN chmod +x ./nef
RUN chmod +x ./entrypoint.sh

CMD ["./entrypoint.sh"]

