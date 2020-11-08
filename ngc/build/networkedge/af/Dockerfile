# SPDX-License-Identifier: Apache-2.0
# Copyright © 2019 Intel Corporation

FROM alpine:3.9

ENV http_proxy=$http_proxy
ENV https_proxy=$https_proxy
ENV no_proxy=localhost,127.0.0.1,afservice,nefservice

RUN apk update && apk add bash
RUN apk add --no-cache \
        wget \
        libc6-compat \
        libstdc++ \
        net-tools \
	      nmap

WORKDIR /root
COPY ./af .
COPY ./entrypoint.sh .

RUN chmod +x ./af
RUN chmod +x ./entrypoint.sh
CMD ["./entrypoint.sh"]

