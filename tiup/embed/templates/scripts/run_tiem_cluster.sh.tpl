#!/bin/bash
set -e

# WARNING: This file was auto-generated. Do not edit!
#          All your edit might be overwritten!
DEPLOY_DIR={{.DeployDir}}

cd "${DEPLOY_DIR}" || exit 1

source ~/.bash_profile

LANG=en_US.UTF-8 \
exec env GODEBUG=madvdontneed=1 bin/cluster-server \
    --host="{{.Host}}" \
    --port="{{.Port}}" \
    --metrics-port="{{.MetricsPort}}" \
    --registry-address="{{.RegistryEndpoints}}" \
    --tracer-address="{{.TracerAddress}}" \
    --deploy-dir="{{.DeployDir}}/bin" \
    --data-dir="{{.DataDir}}" \
    --log-level="{{.LogLevel}}"