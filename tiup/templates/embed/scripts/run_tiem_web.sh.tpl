#!/bin/bash
set -e

# WARNING: This file was auto-generated. Do not edit!
#          All your edit might be overwritten!
DEPLOY_DIR={{.deploy_dir}}

cd "${DEPLOY_DIR}" || exit 1

LANG=en_US.UTF-8 \
exec env GODEBUG=madvdontneed=1 bin/web-server \
    --host="{{.tiem_web_servers.host}}" \
    --port="{{.tiem_web_servers.port}}" \
    --registry-address="{{.registry_endpoints}}" \
    --deploy-dir="{{.tiem_web_servers.deploy_dir}}"
