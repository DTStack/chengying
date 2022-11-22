#!/bin/bash

source /etc/profile

CMD_PATH=`dirname $0`
CMD_HOME=`cd "$CMD_PATH"/./; pwd`

add_crond(){

   CRONDFILF=/var/spool/cron/${USER}
   IFS=$'\n\n'

   for i in `cat crond`
   do
   if [[ "${i}" != "" ]] && [[ -z `sudo grep "${CMD_HOME}" "$CRONDFILF"` ]]; then
   echo ${i} | sudo tee -a ${CRONDFILF} > /dev/null
   fi
   done

}


add_crond