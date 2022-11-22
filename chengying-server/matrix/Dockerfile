#FROM frolvlad/alpine-glibc
#FROM centos:centos7
FROM lgatica/openssh-client

LABEL maintainer="guyan@dtstack.com"

#mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY matrix example-config.yml entrypoint.sh /matrix/

COPY font /matrix/font

RUN sed -i 's|log.dir:.*|log.dir: /tmp/|; s/^\(\s*agent.host\):.*/\1: easy-agent-server:8889/; s/mysqldb.host:.*/mysqldb.host: mysql/; s/mysqldb.user:.*/mysqldb.user: root/; s/mysqldb.password:.*/mysqldb.password: dtstack/; s/api.host:.*/api.host:/; s/api.static-url:.*/api.static-url:/; s|grafana-url:.*|grafana-url: http://grafana:3000|; ' /matrix/example-config.yml && ln -sf /dev/stdout /tmp/matrix.log && ln -sf /dev/stderr /tmp/matrix-error.log && ln -sf /dev/stdout /tmp/matrix-api.log

WORKDIR /matrix

CMD ./entrypoint.sh

EXPOSE 8864
