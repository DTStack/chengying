#!/bin/sh

#############################################################
#  设置环境变量
############################################################
declare serial
declare imageList
PS3="Enter a Number:"
serial=(1 2 3 4)
menu=("Install_Docker" "Load_Image" "Install_EM" "Update")
## Path Info
script_path=$(cd $(dirname $0) && pwd -P)
images_path=${script_path}/images

logo(){
echo -e "
 _____    _    ______   ____  __    _    _   _    _    ____ _____ ____
| ____|  / \  / ___\ \ / /  \/  |  / \  | \ | |  / \  / ___| ____|  _ \
|  _|   / _ \ \___ \\ V /| |\/| | / _ \ |  \| | / _ \| |  _|  _| | |_) |
| |___ / ___ \ ___) || | | |  | |/ ___ \| |\  |/ ___ \ |_| | |___|  _ <
|_____/_/   \_\____/ |_| |_|  |_/_/   \_\_| \_/_/   \_\____|_____|_| \_\
\n
 ___ _   _ ____ _____  _    _     _
|_ _| \ | / ___|_   _|/ \  | |   | |
 | ||  \| \___ \ | | / _ \ | |   | |
 | || |\  |___) || |/ ___ \| |___| |___
|___|_| \_|____/ |_/_/   \_\_____|_____|
\n

DEPRECATED: This script is used to install Easymanager version v4.1.4 with kubernetes v1.16
Copyright (c) 2020 DTStack Inc.
"
echo "======================= Menu ======================="
}

############################################################
#  公共方法
############################################################

Usage(){
echo "Usage: "
echo "  -p      Create kubernetest Local PV"
echo "  -o      Output Easymanager Yaml"
echo
}


check_ip(){
   local IP=$1
   VALID_CHECK=$(echo $IP|awk -F . '$1<=255&&$2<=255&&$3<=255&&$4<=255{print "yes"}')
   if echo $IP|grep -E "^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$" > /dev/null
   then
       if [ $VALID_CHECK == "yes" ]
       then
           return 0
       else
           echo -e "\033[31mIP $IP format error!\033[0m"
           return 1
       fi
    else
        echo -e "\033[31mIP $IP format error!\033[0m"
        return 1
    fi
}

## 创建easymanager
CreateDpSvc(){
    ## 创建命名空间
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/namespace.yaml
    ## 创建secret
    for sc in `ls ${script_path}/emyaml/secret`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/secret/${sc}
    done
    ## 创建rbac
    for rbac in `ls ${script_path}/emyaml/rbac`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig create -f ${script_path}/emyaml/rbac/${rbac}
    done
    ## 创建services对象
    for svc in `ls ${script_path}/emyaml/services`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/services/${svc}
    done
    ## 创建configmap对象
    for cm in `ls ${script_path}/emyaml/configmap`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/configmap/${cm}
    done
    ## 创建PVC资源对象
    for pvcs in  `ls ${script_path}/emyaml/pvc`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/pvc/${pvcs}
    done
    ## 创建deployment资源对象
    for dp in `ls ${script_path}/emyaml/deployment`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/deployment/${dp}
    done
}

## 创建local pv
CreatePvStorage(){
    ## 创建命名空间
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/namespace.yaml
    ## 创建storageclass对象
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/storageclass.yaml
    ## 创建PV资源对象
    for pvs in `ls ${script_path}/emyaml/pv`;do
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/pv/${pvs}
    done
}

Update(){
    ## 更新configmap对象
    for cm in `ls ${script_path}/emyaml/configmap`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/configmap/${cm}
    done
    ## 更新services对象
    for svc in `ls ${script_path}/emyaml/services`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/services/${svc}
    done
    ## 更新deployment资源对象
    for dp in `ls ${script_path}/emyaml/deployment`;do
    sleep 1
    ${script_path}/kubeconfig/kubectl --kubeconfig=${script_path}/kubeconfig/kubeconfig apply -f ${script_path}/emyaml/deployment/${dp}
    done
}

####################################################
# 安装docker，导入镜像以及上传镜像
####################################################

Install_Docker(){
    echo -e '\n-----------------------------------------------------------------------------'
    echo -e "\033[32mINFO: Install Docker ! \033[0m"
    echo '------------------------------------------------------------------------------'
    which docker
    if [ $? -ne 0 ]; then
    sh ${script_path}/install_docker/install_docker.sh ${script_path}/install_docker/docker-18.06.1-ce.tgz
    else
       echo -e "\033[31mdocker runtime already exists !\033[0m"
    fi
}

Load_Image(){
     if [ $(ps -C dockerd --no-heading | wc -l) -eq 0 ]
        then
        echo -e '\n-----------------------------------------------------------------------------'
        echo -e "\033[31mERROR: Is the docker daemon running? Please check \033[0m"
        echo '------------------------------------------------------------------------------'
        return 1
    else
        for i in `ls ${images_path}`
        do
            images=$(docker load -i ${images_path}/${i}|grep -i "loaded image:" |awk '{print $3}')
            k=$((index++))
            imageList[${k}]=${images}
        done
        return 0
    fi
}

Push_Image(){
    echo -e '\n-----------------------------------------------------------------------------'
    echo -e "\033[32mINFO: Load Image,Please provide host,username,password information! \033[0m"
    echo '------------------------------------------------------------------------------'
    read -p "[Input Repository(example 10.10.10.10/test)]:" repository
    if [ -z ${repository} ];then
        echo -e "\033[31mInput Null ! Please select the module again -_-!\033[0m"
        continue
    fi
    ## ip地址格式检查
    formatip=${repository%%/*}
    ip=${formatip%:*}
    check_ip $ip
    [ $? -eq 1 ] && continue

    read -p "[Input Repository Username]:" repositoryUser
    read -p "[Input RegRepositoryistry Password]:" repositoryPass

    if [ -z ${repositoryUser} ];then
        echo -e "\033[31mInput Null ! Please select the module again -_-!\033[0m"
        continue
    fi

    if [ -z ${repositoryPass} ];then
        echo -e "\033[31mInput Null ! Please select the module again -_-!\033[0m"
        continue
    fi

    ## 登录仓库
    docker login ${formatip} -u${repositoryUser} -p${repositoryPass}
    [ $? -eq 1 ] && echo -e "\033[31mdocker login failed!\033[0m" && continue

    ## 导入镜像到本地,打标签,推送
    echo "----------------------- Load image And Push image -----------------------"
    Load_Image
    if [ $? -eq 0 ]
    then
        ## 打标签he推送镜像
        for img in ${imageList[@]};do
            docker tag ${img} ${repository}/${img##*/}
            docker push ${repository}/${img##*/}
        done
    else
       echo -e "\033[31m Load image failed ! -_-!\033[0m"
    fi
}

################################################
# 使用token或者kubeconfig部署em
################################################
Install_emORpv()
{
    echo -e '\n-----------------------------------------------------------------------------'
    echo -e "\033[32mINFOR: Install Easymanager with kubeconfig or token \033[0m"
    echo '------------------------------------------------------------------------------'
    read -p '[Select authentication mode "kubeconfig" or "token"]:' sec
    if [ -z ${sec} ];then
        echo -e "\033[31mInput Null ! Please select the module again -_-!\033[0m"
        exit 1
    fi

    case ${sec} in
    "token")
        if [ -e ${script_path}/kubeconfig/kubeconfig -a -s ${script_path}/kubeconfig/kubeconfig ]
        then
            ## 创建k8s资源对象
            $1
            exit 1
        else
            ## 使用token创建k8s资源对象
            read -p "[Input kubernetes Api(example 10.10.10.10:6443)]:" k8sapi
            if [ -z ${k8sapi} ];then
            echo -e "\033[31mInput Null ! Please Attention -_-!\033[0m"
            exit 1
            fi
            ip=${k8sapi%:*}
            check_ip $ip
            [ $? -eq 1 ] && exit 1

            ## 使用token创建k8s资源对象
            read -p "[Input Token]:" apitoken
            if [ -z ${apitoken} ];then
                echo -e "\033[31mInput Null ! Please select the module again -_-!\033[0m"
                exit 1
            fi
            ## 解析生成kubeconfig文件
            ${script_path}/kubeconfig/kubectl config set-cluster userk8s \
            --server=https://${k8sapi} \
            --kubeconfig=${script_path}/kubeconfig/kubeconfig \
            --insecure-skip-tls-verify

            ${script_path}/kubeconfig/kubectl config set-credentials tkuser \
            --token="${apitoken}" \
            --kubeconfig=${script_path}/kubeconfig/kubeconfig

            ${script_path}/kubeconfig/kubectl config set-context tkuser@userk8s \
            --cluster=userk8s \
            --user=tkuser \
            --kubeconfig=${script_path}/kubeconfig/kubeconfig

            ${script_path}/kubeconfig/kubectl config use-context tkuser@userk8s \
            --kubeconfig=${script_path}/kubeconfig/kubeconfig
            ## 创建k8s资源对象
            $1
            exit 1
        fi
    break
    ;;
    "kubeconfig")
        ## 检查kubeconfig
        if [ -e ${script_path}/kubeconfig/kubeconfig -a -s ${script_path}/kubeconfig/kubeconfig ]
        then
            ## 创建k8s资源对象
            $1
            exit 1
        else
            echo -e "\033[31m${script_path}/kubeconfig/kubeconfig not exist or is empty !\033[0m"
            exit 1
        fi
    ;;
    *)
      echo -e "\033[31mPlease Select authentication mode 'kubeconfig' or 'token'\033[0m"
      exit 1
    ;;
    esac
}

#################
# 菜单
#################
menu(){
select var in ${menu[@]}
do
   ######## Serial Number Range
   if ! echo ${serial[@]} | grep -q $REPLY;then
   echo "Error Number,Please Enter[1-4]"
   continue
   fi
   case $var in
   Install_Docker)
      Install_Docker
      break
   ;;
   Load_Image)
      Push_Image
      break
   ;;
   Install_EM)
      Install_emORpv CreateDpSvc
      break
   ;;
   Update)
      Push_Image
      ## 更新deployment资源对象
      Update
      break
   ;;
   *)
      break
   ;;
   esac

done
}

## 参数解析
while getopts "hpo" arg
do
    case $arg in
    h)
        Usage
        exit 1
        ;;
    p)
        sh ${script_path}/translate.sh
        Install_emORpv CreatePvStorage
        exit 1
        ;;
    o)
        sh ${script_path}/outputyaml.sh
        exit 1
        ;;
    *)
        Usage
        exit 1
        ;;
    esac
done

## 默认执行
sh ${script_path}/translate.sh
logo
menu