#!/bin/bash
set -x
path=/opt/dtstack/kerberos/kdcserver

ips="{{.Join "zookeeper_ip" " "}}"
hosts="{{.JoinHost "zookeeper_ip" " "}}"

zookeeper_host="{{.JoinHost "zookeeper_ip" " "}}"

ssh -t {{.kdcserver_ip}} " $path/create_principal_1213.sh zookeeper $hosts $ips "
exit

