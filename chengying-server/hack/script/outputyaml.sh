#!/bin/sh

declare -A dict
script_path=$(cd $(dirname $0) && pwd -P)
/bin/sh ${script_path}/translate.sh

## 输出namespace
cat ${script_path}/emyaml/namespace.yaml >> ${script_path}/easymanager-all.yaml
echo "---" >> ${script_path}/easymanager-all.yaml

## 输出storageclass
cat ${script_path}/emyaml/storageclass.yaml >> ${script_path}/easymanager-all.yaml
echo "---" >> ${script_path}/easymanager-all.yaml

## 输出secret
for sc in `ls ${script_path}/emyaml/secret`;do
    cat ${script_path}/emyaml/secret/${sc} >> ${script_path}/easymanager-all.yaml
    echo "---" >> ${script_path}/easymanager-all.yaml
done

## 输出rbac
for sc in `ls ${script_path}/emyaml/rbac`;do
    cat ${script_path}/emyaml/rbac/${rbac} >> ${script_path}/easymanager-all.yaml
    echo "---" >> ${script_path}/easymanager-all.yaml
done

## 输出deployment
for dp in `ls ${script_path}/emyaml/deployment`;do
    cat ${script_path}/emyaml/deployment/${dp} >> ${script_path}/easymanager-all.yaml
    echo "---" >> ${script_path}/easymanager-all.yaml
done

## 输出services
for svc in `ls ${script_path}/emyaml/services`;do
    cat ${script_path}/emyaml/services/${svc} >> ${script_path}/easymanager-all.yaml
    echo "---" >> ${script_path}/easymanager-all.yaml
done

## 输出pvc
for pvcs in `ls ${script_path}/emyaml/pvc`;do
    cat ${script_path}/emyaml/pvc/${pvcs} >> ${script_path}/easymanager-all.yaml
    echo "---" >> ${script_path}/easymanager-all.yaml
done

## 输出pv
for pvs in `ls ${script_path}/emyaml/pv`;do
    cat ${script_path}/emyaml/pv/${pvs} >> ${script_path}/easymanager-all.yaml
    echo "---" >> ${script_path}/easymanager-all.yaml
done

## 输出configmap
for cm in `ls ${script_path}/emyaml/configmap`;do
    cat ${script_path}/emyaml/configmap/${cm} >> ${script_path}/easymanager-all.yaml
    echo "---" >> ${script_path}/easymanager-all.yaml
done