#!/bin/sh
set -e

: ${ENVIRONMENT:=missing}
: ${SUXEND_PORT:=8080}
: ${SUXEND_HOST:=0.0.0.0}
: ${SUXEND_STATIC:=/dist}
: ${SUXEND_LOG_LEVEL:=info}
: ${SUXEND_LOG_ENVIRONMENT:=production}
: ${SUXEND_NEXUS_ADDRESS:=nexus.example.com}
: ${SUXEND_NEXUS_SVC_ADDRESS:=nexus.example.com}
: ${SUXEND_NEXUS_REGISTRY_ADDRESS:=container.example.com}


if [ "$1" = 'suxend' ]; then

exec suxend -env=${ENVIRONMENT} \
    -host=${SUXEND_HOST} \
	-port=${SUXEND_PORT} \
	-static=${SUXEND_STATIC} \
	-gcloud.project=${SUXEND_GCLOUD_PROJECT} \
	-log.environment=${SUXEND_LOG_ENVIRONMENT} \
	-log.level=${SUXEND_LOG_LEVEL} \
	-nexus.address=${SUXEND_NEXUS_ADDRESS} \
	-nexus.svc.address=${SUXEND_NEXUS_SVC_ADDRESS} \
	-nexus.svc.authToken=${SUXEND_NEXUS_SVC_AUTH_TOKEN} \
	-nexus.repository=${SUXEND_NEXUS_REPOSITORY} \
	-nexus.registry.address=${SUXEND_NEXUS_REGISTRY_ADDRESS}
fi

exec "$@"