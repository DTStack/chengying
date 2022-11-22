#!/bin/sh
# env
script_path=$(cd $(dirname $0) && pwd -P)
yaml_path=${script_path}/emyaml

if [ -e ${script_path}/kubeconfig/kubeconfig -a -s ${script_path}/kubeconfig/kubeconfig ]
then
## 删除命名空间
#${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${yaml_path}/namespace.yaml

    ## 删除deployment资源对象
    for dp in `ls ${script_path}/emyaml/deployment`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${script_path}/emyaml/deployment/${dp}
    done

    ## 删除services对象
    for svc in `ls ${script_path}/emyaml/services`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${script_path}/emyaml/services/${svc}
    done


    ## 删除configmap对象
    for cm in `ls ${script_path}/emyaml/configmap`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${script_path}/emyaml/configmap/${cm}
    done

    sleep 10

    ## 删除storageclass对象
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${script_path}/emyaml/storageclass.yaml

    ## 删除PVC资源对象
    for pvcs in  `ls ${script_path}/emyaml/pvc`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${script_path}/emyaml/pvc/${pvcs}
    done

    ## 创建PV资源对象
    for pvs in `ls ${script_path}/emyaml/pv`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${script_path}/emyaml/pv/${pvs}
    done

    ## 删除secret
    for sc in `ls ${script_path}/emyaml/secret`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${script_path}/emyaml/secret/${sc}
    done

    ## 删除rbac
    for rbac in `ls ${script_path}/emyaml/rbac`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig delete -f ${script_path}/emyaml/rbac/${rbac}
    done

else
    echo -e "\033[31m${script_path}/kubeconfig/kubeconfig not exist or is empty !\033[0m"
#    return 1
fi
