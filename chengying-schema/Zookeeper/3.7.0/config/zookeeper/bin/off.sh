#!/bin/bash
set -x
ssh -t {{.kdcserver_ip}}  "rm -f {{.zookeeper_keytab}}"

ssh -t {{.kdcserver_ip}}  "rm -f {{.zkcli_keytab}}"