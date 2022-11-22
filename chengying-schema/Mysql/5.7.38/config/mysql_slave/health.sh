#!/bin/bash

#usage ./port_status.sh localhost 80 8080 8088
#parameter 1 is host
#parameters after 1 are list of ports to check
#all ports are ok echo 1, else echo 0 


if [ $# -lt 2 ] ; then
  echo "USAGE: $0 host [ports]" 
  exit
fi

./bin/mysqlshow -h$1 -u$2 -p$3
