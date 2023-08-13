#!/bin/sh

#set -e
#set -u


if [ ! -d  "/matrix/easyagent/DTBase" ]; then
 mv /matrix/tmp-easyagent/* /matrix/easyagent
fi

DTSTACK_HOME="/matrix"
DTSTACK_INIT="${DTSTACK_HOME}/init"

if [ "$(ls ${DTSTACK_INIT})" ]; then
  for init in ${DTSTACK_INIT}/*.sh; do
    . $init
  done
fi

#nginx -g "user admin;daemon on;"
#nginx -g "daemon on;"
./matrix -c example-config.yml --debug
sleep 5
tail -f /dev/null

