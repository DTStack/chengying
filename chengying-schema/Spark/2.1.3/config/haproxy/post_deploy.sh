#!/bin/bash

set -e

USER=`whoami`
sed -i  s/admin/$USER/g conf/haproxy.cfg