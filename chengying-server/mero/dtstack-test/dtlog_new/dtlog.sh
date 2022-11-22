#!/bin/bash

echo "dtlog $@"
i=0
while true
do
    echo dtlog $i "$DTLOG:$ABC" version `cat current/version.txt` >> $2
    i=$[i+1]
    sleep 5
done
