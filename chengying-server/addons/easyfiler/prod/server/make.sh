#!/bin/sh

gox -os=linux -arch=amd64
mv server_linux_amd64 easyfiler
cp -f easyfiler easyfiler-em/easyfiler/sbin/

./mero ./easyfiler-em
