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

