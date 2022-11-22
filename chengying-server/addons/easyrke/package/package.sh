#!/bin/sh


if [ $# -lt 1 ] ; then
  echo "USAGE: $0 k8sversion"
  exit
fi

k8sversion=$1

rm -rf DTK8S-without-images*
#确保Jenkins静态资源存在
wget http://172.16.10.109:88/packages/easymanager/DTK8S-without-images.tgz

tar -xvf DTK8S-without-images.tgz
cp ../images/* DTK8S-without-images/docker/images/

sed -i "" "s|^\(\s*kubernetes_version\):.*|\1: $k8sversion|" "DTK8S-without-images/rke/config/cluster.yml"
sed -i "" "s|^\(\s*product_version\):.*|\1: DTK8S-$k8sversion|" "DTK8S-without-images/schema.yml"

#确保mero可执行文件存在系统PATH
../mero ./DTK8S-without-images
