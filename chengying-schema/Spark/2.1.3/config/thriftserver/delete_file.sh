#! /bin/bash
find /tmp/admin/ -name '*.pipeout'  -amin +360 -ls  -exec rm {} \; > delete_file.log

