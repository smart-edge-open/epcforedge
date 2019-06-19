```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation and Smart-Edge.com, Inc.
```

# http

This is the guide how to bringup  nginx as HTTPS server with reference configuration.


## Overview

- Run install_nginx.sh to install nginx and fcgi.
- After installation, use below reference configuration to replace the defaut configuration: /etc/nginx/nginx.conf 
  Examples:
```text
   worker_processes 1;
   worker_cpu_affinity 0100000000;
   error_log  /var/log/nginx.log  debug;
   pid        /var/log/nginx.pid;
   events {
      worker_connections  1024;
   }

   http {
      include                         mime.types;
      default_type                    application/octet-stream;
      sendfile                        on;
      keepalive_timeout               65;
      server {
        listen       8080 ssl;
        ssl_certificate /etc/nginx/ssl/epc.crt;
        ssl_certificate_key /etc/nginx/ssl/epc.key;
        server_name  localhost;
        location /userplanes {
                fastcgi_pass  127.0.0.1:9999;
                include       fastcgi_params;
                fastcgi_param HTTPS on;
        }
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
      }
   }
```

-  To generate the certificate and key, use below commands:
```text
   mkdir /etc/nginx/ssl/
   sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout /etc/nginx/ssl/epc.key -out /etc/nginx/ssl/epc.crt -subj "/C=US/ST=Epc/L=Epc/O=Epc/OU=Epc/CN=epc.oam"
```

## Run Guide
- Type nginx directly to run it.
- If you want to change /etc/nginx/nginx.conf , type: nginx-t  and then nginx -s reload. Thus nginx will reload configuration and re-start.



