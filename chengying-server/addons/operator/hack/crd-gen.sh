#!/usr/bin/env bash

## find controller-tools pkg
get_controllergen_pkg(){
    local IFS=" "
    s=$(grep -i sigs.k8s.io/controller-tools $(dirname "${BASH_SOURCE[0]}")/../../../go.mod | sed -e 's/^[[:space:]]*//g' | sed -e 's/[[:space:]]*$//g')
    array=($s)
    CONTROLLER_GEN_PKG=$(go env GOPATH)/pkg/mod/${array[0]}@${array[1]}
    echo $CONTROLLER_GEN_PKG
}

CONTROLLER_GEN_PKG=$(get_controllergen_pkg)

## find go bin path
GOBIN="$(go env GOBIN)"
gobin="${GOBIN:-$(go env GOPATH)/bin}"


## build controller-gen
if [ ! -e ${gobin}/controller-gen ];then
    cd $CONTROLLER_GEN_PKG
    go install ./cmd/controller-gen
fi

echo "Generating crd yaml file in:$(dirname "${BASH_SOURCE[0]}")/../deploy/crds"
"${gobin}/controller-gen" "crd:trivialVersions=true,preserveUnknownFields=false" paths="./..." output:crd:artifacts:config=$(dirname "${BASH_SOURCE[0]}")/../deploy/crds

