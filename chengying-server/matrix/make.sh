#!/bin/sh

gox -os=linux -arch=amd64
mv matrix_linux_amd64 matrix
docker build -t 172.16.8.120:5443/dtstack-dev/matrix:4.1.9-workloadAdaptBigdata-workloadSvcName .
docker push 172.16.8.120:5443/dtstack-dev/matrix:4.1.9-workloadAdaptBigdata-workloadSvcName

