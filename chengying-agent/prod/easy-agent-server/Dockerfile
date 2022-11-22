FROM frolvlad/alpine-glibc:alpine-3.7

RUN apk --update add netcat-openbsd

COPY easy-agent-server example-config.yml entrypoint.sh /easy-agent-server/

RUN sed -i 's|log.dir:.*|log.dir: /tmp/|; s/^\(\s*host\):.*/\1: "matrix:8864"/; s/mysqldb.host:.*/mysqldb.host: mysql/; s/mysqldb.user:.*/mysqldb.user: root/; s/mysqldb.password:.*/mysqldb.password: dtstack/' /easy-agent-server/example-config.yml && ln -sf /dev/stdout /tmp/easyagent-server.log && ln -sf /dev/stderr /tmp/easyagent-server-error.log && ln -sf /dev/stdout /tmp/api.log

WORKDIR /easy-agent-server
COPY bin /easy-agent-server/

CMD ./entrypoint.sh

EXPOSE 8889
EXPOSE 8890
