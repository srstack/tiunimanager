worker_processes  1;

events {
    worker_connections  1024;
}

http {
    include       mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  {{.LogDir}}/access.log  main;

    upstream openapi-servers {
        least_conn;

        upsync {{.IP}}:{{.Port}}/etcd/v2/keys/micro/registry/openapi-server upsync_timeout=6m upsync_interval=500ms upsync_type=etcd strong_dependency=off;
        upsync_dump_path {{.DeployDir}}/conf/server_list.conf;
        include {{.DeployDir}}/conf/server_list.conf;
    }

    upstream etcdcluster {
{{- range $idx, $addr := .RegistryEndpoints}}
        server {{$addr}} weight=1 fail_timeout=10 max_fails=3;
{{- end}}
    }

    gzip  on;

    server {
        listen {{.Port}};
        server_name  {{.ServerName}};

        location / {
            root   html;
            index  index.html index.htm;
            try_files $uri $uri/ /index.html;
        }

{{- if not .EnableHttps }}
        location ^~ /api {
            proxy_pass http://openapi-servers;
        }

        location ~ ^/(swagger|system|web)/ {
            proxy_pass http://openapi-servers;
        }
{{- end}}

        location ^~/etcd/ {
            proxy_pass http://etcdcluster/;
        }

        location ~ ^/env {
            default_type application/json;
            return 200 '{"protocol": "{{.Protocol}}", "tlsPort": {{.TlsPort}}, "service": {"grafana": "http://{{.GrafanaAddress}}/d/tiem000001/tiem-server?orgId=1&refresh=10s&kiosk=tv", "kibana": "http://{{.KibanaAddress}}/app/discover", "alert": "http://{{.AlertManagerAddress}}", "tracer": "http://{{.TracerAddress}}"}}';
        }

        location = /upstream_show {
            upstream_show;
        }
    }

{{- if .EnableHttps }}

    server {
        listen {{.TlsPort}} ssl;
        server_name  {{.ServerName}};

        ssl_certificate  {{.DeployDir}}/cert/server.crt;
        ssl_certificate_key {{.DeployDir}}/cert/server.key;
        server_tokens off;

        fastcgi_param   HTTPS               on;
        fastcgi_param   HTTP_SCHEME         https;

        location ^~ /api {
            proxy_pass https://openapi-servers;
        }

        location ~ ^/(swagger|system|web)/ {
            proxy_pass https://openapi-servers;
        }
    }
{{- end}}
}