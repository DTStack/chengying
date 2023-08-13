#!/bin/sh

#set -e
#set -u

DTSTACK_HOME="/easy-agent-server"
DTSTACK_INIT="${DTSTACK_HOME}/init"

if [ "$(ls ${DTSTACK_INIT})" ]; then
  for init in ${DTSTACK_INIT}/*.sh; do
    . $init
  done
fi

#nginx -g "user admin;daemon on;"
#nginx -g "daemon on;"
./easy-agent-server -c example-config.yml --debug
sleep 5
tail -f /dev/null
