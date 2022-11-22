#!/bin/bash

set -e

cur=$(pwd)
case "$(uname)" in
Linux)
	bin_path=$(readlink -f $(dirname $0))
	;;
*)
	bin_path=$(
		cd $(dirname $0)
		pwd
	)
	;;
esac

HADOOP_HOME=/data/hadoop_base
USER=$(whoami)
cp -rf miniconda* /data/
check_kernel() {

	cpu_kernel=$(uname -a | awk '{print $3}')

	if [ "$cpu_kernel" = "3.10.0-327.el7.x86_64" ]; then
		sed -i s/"org.apache.hadoop.yarn.server.nodemanager.util.CgroupsLCEResourcesHandler"/"org.apache.hadoop.yarn.server.nodemanager.util.DefaultLCEResourcesHandler"/g $HADOOP_HOME/etc/hadoop/yarn-site.xml
		sed -i s/"org.apache.hadoop.yarn.server.nodemanager.LinuxContainerExecutor"/"org.apache.hadoop.yarn.server.nodemanager.DefaultContainerExecutor"/g $HADOOP_HOME/etc/hadoop/yarn-site.xml
	fi
}

CGROUP_HOME=/sys/fs/cgroup
cgroup_config_path=/etc/yarn-cgroup

check_cgroup() {
	if [ -d $CGROUP_HOME ]; then
		sudo mkdir -p /sys/fs/cgroup/cpu,cpuacct/hadoop-yarn/
		sudo chown -R $USER:$USER /sys/fs/cgroup/cpu,cpuacct/hadoop-yarn

		##chmod to container-executor
		sudo chown root:$USER $HADOOP_HOME/bin/container-executor
		sudo chmod 6050 $HADOOP_HOME/bin/container-executor
	fi
}

check_cgroup_config() {
	if [ ! -d $cgroup_config_path ]; then
		sudo mkdir -p $cgroup_config_path
		sudo cp -rp $HADOOP_HOME/etc/hadoop/container-executor.cfg $cgroup_config_path/container-executor.cfg
		sudo chown root:$USER -R $cgroup_config_path
	fi
}

check_kernel
check_cgroup
check_cgroup_config
