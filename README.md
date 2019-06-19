```text
Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

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
