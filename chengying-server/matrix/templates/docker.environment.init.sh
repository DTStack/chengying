#!/bin/bash

STATIC_HOST='{{.STATIC_HOST}}'
REGISTRY_HOST='{{.MATRIX_IP}}'

install_docker(){
    curl -L -O -s "$STATIC_HOST/easyagent/dtstack-docker.tar.gz"
    tar -xf "dtstack-docker.tar.gz" -C /opt/dtstack/

    sh /opt/dtstack/dtstack-docker/install.sh $REGISTRY_HOST
}

install_docker