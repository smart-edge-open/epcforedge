```text
SPDX-License-Identifier: Apache-2.0
Copyright © 2019 Intel Corporation.
```

# http

This is the guide how to bringup  nginx as HTTP server with reference configuration.


## Overview

- Run install_nginx.sh to install nginx and fcgi.
- After installation, use below reference configuration to replace the defaut configuration: /etc/nginx/nginx.conf 
- The nginx configuration should be based on HTTP and CORS enablement.
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
    #tcp_nopush                     on;
    keepalive_timeout               65;
    #gzip  on;

    server {
        listen       8080;
        server_name  localhost;
        #charset     utf-8;
        add_header "Access-Control-Allow-Origin" "*" always;
        add_header "Access-Control-Allow-Headers" "Content-Type" always;

        location /userplanes {
                if ($request_method = "OPTIONS") {
                   add_header "Access-Control-Allow-Origin" "*" always;
                   add_header "Access-Control-Allow-Headers" "Content-Type" always;
                   add_header "Access-Control-Allow-Methods" "GET,POST,PATCH,DELETE" always;
                   return 200;
                }

                if ($request_method = "GET") {
                  add_header "Access-Control-Allow-Origin" "*" always;
                  add_header "Content-Type" "application/json" always;
                }

                if ($request_method = "DELETE") {
                  add_header "Access-Control-Allow-Origin" "*" always;
                  add_header "Content-Type" "application/json" always;
                }

                if ($request_method = "PATCH") {
                  add_header "Access-Control-Allow-Origin" "*" always;
                  add_header "Content-Type" "application/json" always;
                }

                if ($request_method = "POST") {
                  add_header "Access-Control-Allow-Origin" "*" always;
                  add_header "Content-Type" "application/json" always;
                }

                fastcgi_pass  127.0.0.1:9999;
                include       fastcgi_params;
                fastcgi_param HTTPS off;
        }

        #error_page  404              /404.html;
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }
}

```


## Run Guide
- Type nginx directly to run it.
- If you want to change /etc/nginx/nginx.conf , type: nginx-t  and then nginx -s reload. Thus nginx will reload configuration and re-start.
	    
