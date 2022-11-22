#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

if [ "$#" -lt 2 ] || [ "${1}" == "--help" ]; then
  cat <<EOF
Usage: $(basename "$0") <apis-package> <groups-versions>

  <apis-package>      the external types dir (e.g. github.com/example/api or github.com/example/project/pkg/apis).
  <groups-versions>   the groups and their versions in the format "groupA:v1,v2 groupB:v1 groupC:v2", relative
                      to <api-package>.

Examples:
  $(basename "$0") github.com/example/project/pkg/apis "foo:v1 bar:v1alpha1,v1beta1"
EOF
  exit 0
fi


APIS_PKG="$1"
GROUPS_WITH_VERSIONS="$2"
shift 2

get_codegen_pkg(){
    local IFS=" "
    s=$(grep -i k8s.io/code-generator $(dirname "${BASH_SOURCE[0]}")/../../../go.mod | sed -e 's/^[[:space:]]*//g' | sed -e 's/[[:space:]]*$//g')
    array=($s)
    CODEGEN_PKG=$(go env GOPATH)/pkg/mod/${array[0]}@${array[1]}
    echo $CODEGEN_PKG
}

codegen_join(){
    local IFS="$1"
    shift
    echo "$*"
}

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=$(get_codegen_pkg)

FQ_APIS=() # e.g. k8s.io/api/apps/v1
for GVs in ${GROUPS_WITH_VERSIONS}; do
  IFS=: read -r G Vs <<<"${GVs}"

  # enumerate versions
  for V in ${Vs//,/ }; do
    FQ_APIS+=("${APIS_PKG}/${G}/${V}")
  done
done

GOBIN="$(go env GOBIN)"
gobin="${GOBIN:-$(go env GOPATH)/bin}"

if [ ! -e ${gobin}/deepcopy-gen ];then
    cd $CODEGEN_PKG
    go install ./cmd/deepcopy-gen
fi

echo "Generating deepcopy funcs"

"${gobin}/deepcopy-gen" --input-dirs "$(codegen_join , "${FQ_APIS[@]}")" -O zz_generated.deepcopy --bounding-dirs "${APIS_PKG}" --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

for fq_api in ${FQ_APIS[@]}; do
  echo $fq_api
  mv $(go env GOPATH)/src/$fq_api/zz_generated.deepcopy.go ${SCRIPT_ROOT}/${fq_api#*operator/}
done
