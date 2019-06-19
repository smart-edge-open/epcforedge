```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation and Smart-Edge.com, Inc.
```
# Purpose

This document is intended for EPC OAMAgent setup and serves as a guide on setting up agent that will play as interface between controller and EPC control plane,

# EPC OAMAgent Community Edition

This repository contains the project for the EPC OAMAgent Community Edition.

## Directory structure

The EPC OAMAgent supports EPC-Cplane configuration from MEC controller. So provides two basic components as below:

- backend: CGI backend source files to handle HTTP(S) commmunciation with controller
- http: only provide guide about how to configure and start NGINX , also including how to generate self-signed certification files for HTTPS

## Installation, Build and Run

- Setup, Generate Self-signed certification files and Run NGINX according to readme in sub-directory: http
- Setup, Build and Run HTTPS based backend according to readme in the sub-directory: backend
- Notice: It is running on CentOS 7.6 x86_64 operating system.

# Troubleshooting

* Script stops/freezes at fetching packages

  Make sure that proxy is configured properly in operating system.

* Controller not able to communicate with OAMAgent

  Make sure that it is not caused by HTTPS proxy setting for the operating system that OAMAgent is running on.
  If finding "Permission denied" in the nginx log, can use command: setenforce=0

* Log files

  EPC OAMAgen uses syslog as logging tool. So can find debug information from /var/log/message.


