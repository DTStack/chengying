#!/bin/sh
script_path=$(cd $(dirname $0) && pwd -P)

declare -A dict

## prometheus等配置文件转换为yaml文件
configToyaml(){
for parent in `ls ${script_path}/conftmp/configfile`
do
local ifstp=$IFS
for file in `ls ${script_path}/conftmp/configfile/${parent}`
do
    local IFS=""
    while read line;do
        filevalue+="$line\n"
    done < ${script_path}/conftmp/configfile/${parent}/${file}
    filename=${file%%.*}
echo "
apiVersion: v1
kind: ConfigMap
metadata:
  name: ${parent}
  namespace: {{.NameSpace}}
data:
 ${file}: \"${filevalue}\"" > ${script_path}/conftmp/emyamltmp/configmap/em-${parent}-cm.yaml

local IFS=$ifstp
local filevalue=
done

done
}

configToyaml

## 根据配置模版解析生成资源yaml对象文件
if [ -e ${script_path}/conftmp/config.tmp -a -s ${script_path}/conftmp/config.tmp ]
then
    grep -v -E "\[|^$" ${script_path}/conftmp/config.tmp > ${script_path}/conftmp/tmp.1
    while read line
    do
        k=${line%=*}
        v=${line#*=}
        #echo "$k:$v"
        dict["${k}"]=${v}
    done <${script_path}/conftmp/tmp.1

    for key in ${!dict[*]}
    do
        st+="s/{{.$key}}/${dict[$key]}/g;"
    done

    ## 解析deployment模版生成实际对象
    for dp in `ls ${script_path}/conftmp/emyamltmp/deployment`;do
        sed -e "${st}" ${script_path}/conftmp/emyamltmp/deployment/${dp} > ${script_path}/emyaml/deployment/${dp}
    done

    ## 解析service模版生成实际对象
    for svc in `ls ${script_path}/conftmp/emyamltmp/services`;do
        sed -e "${st}" ${script_path}/conftmp/emyamltmp/services/${svc} > ${script_path}/emyaml/services/${svc}
    done

    ## 解析pvc模版生成实际对象
    for pvcs in `ls ${script_path}/conftmp/emyamltmp/pvc`;do
        sed -e "${st}" ${script_path}/conftmp/emyamltmp/pvc/${pvcs} > ${script_path}/emyaml/pvc/${pvcs}
    done

    ## 解析pv模版生成实际对象
    for pvs in `ls ${script_path}/conftmp/emyamltmp/pv`;do
        sed -e "${st}" ${script_path}/conftmp/emyamltmp/pv/${pvs} > ${script_path}/emyaml/pv/${pvs}
    done

    ## 解析configmap模版生成实际对象
    for cm in `ls ${script_path}/conftmp/emyamltmp/configmap`;do
        sed -e "${st}" ${script_path}/conftmp/emyamltmp/configmap/${cm} > ${script_path}/emyaml/configmap/${cm}
    done

    ## 解析secret模版生成实际对象
    for sc in `ls ${script_path}/conftmp/emyamltmp/secret`;do
        sed -e "${st}" ${script_path}/conftmp/emyamltmp/secret/${sc} > ${script_path}/emyaml/secret/${sc}
    done

    ## 解析rbac模版生成实际对象
    for rbac in `ls ${script_path}/conftmp/emyamltmp/rbac`;do
        sed -e "${st}" ${script_path}/conftmp/emyamltmp/rbac/${rbac} > ${script_path}/emyaml/rbac/${rbac}
    done

    ## 解析storageclass模版生成实际对象
    sed -e "${st}" ${script_path}/conftmp/emyamltmp/storageclass.yaml > ${script_path}/emyaml/storageclass.yaml

    ## 解析namespace模版生成实际对象
    sed -e "${st}" ${script_path}/conftmp/emyamltmp/namespace.yaml > ${script_path}/emyaml/namespace.yaml
else
    echo -e "\033[31m${script_path}/conftmp/config.tmp not exist or is empty !\033[0m"
fi