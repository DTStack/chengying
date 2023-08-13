#!/bin/bash

set -e
set -u

DTSTACK_HOME="/dtstack"
DTSTACK_INIT="${DTSTACK_HOME}/init"

if [ "$(ls ${DTSTACK_INIT})" ]; then
  for init in ${DTSTACK_INIT}/*.sh; do
    . $init
  done
fi

## clean up pid files
rm -f /tmp/*.pid
rm -f /var/run/*.pid

#################### app start up
APP_HOME=/dt-alert
rm -f ${APP_HOME}/run/*.pid
## chown -R admin:admin ${APP_HOME}

sleep 5

