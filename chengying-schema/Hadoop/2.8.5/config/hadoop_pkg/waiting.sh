#!/bin/bash

USER=`whoami`
sed -i s/\${user\\.name}/$USER/g etc/hadoop/core-site.xml

tail -f -
