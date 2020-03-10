/*******************************************************************************
* Copyright 2019 Intel Corporation. All rights reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*******************************************************************************/
#ifndef __OAMAGENT__LOCALCFG__
#define __OAMAGENT__LOCALCFG__

/* ------------------------------------------------------------------------- */
/* Dependencies */
/* ------------------------------------------------------------------------- */
#include <stdio.h>
#include <string>

using namespace std;

/* ------------------------------------------------------------------------- */
/* Constants */
/* ------------------------------------------------------------------------- */

extern string localcfg_pgw_ipaddress;
extern string localcfg_pgw_port;
extern string localcfg_sgw_ipaddress;
extern string localcfg_sgw_port;

extern string localcfg_nginx_hostip;
extern string localcfg_nginx_port;
extern string localcfg_nginx_fcgipass;

/* ------------------------------------------------------------------------- */
/* Public Functions */
/* ------------------------------------------------------------------------- */
int oamReadCfgJsonFile(void);

#endif

