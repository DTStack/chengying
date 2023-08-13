FROM registry.docker-cn.com/frolvlad/alpine-glibc:alpine-3.7

LABEL maintainer="guyan@dtstack.com"

ENV \
  OSS_EXTERNAL=dtstack-download.oss-cn-hangzhou.aliyuncs.com \
  OSS_INTERNAL=dtstack-download.oss-cn-hangzhou-internal.aliyuncs.com
  
#RUN apk --update add curl netcat-openbsd && mkdir /easy-agent-server
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk --update add curl netcat-openbsd && mkdir /easy-agent-server

# add easy-agent-server
RUN cd /tmp \
  && curl -L http://${OSS_EXTERNAL}/insight/em/2.5/easyagent/easyagent-2.5.tgz > easyagent.tgz \
  && tar -zxf easyagent.tgz && mv bin/easy-agent-server /easy-agent-server \
  && rm -f easyagent.tgz && rm -rf bin

COPY drpeco/example-config.yml /easy-agent-server/
COPY drpeco/dtstack/init /easy-agent-server/init
COPY drpeco/dtstack/bootstrap.sh /easy-agent-server/
#RUN sed -i 's|log.dir:.*|log.dir: /tmp/|; s/^\(\s*host\):.*/\1: "matrix:8864"/; s/mysqldb.host:.*/mysqldb.host: mysql/; s/mysqldb.user:.*/mysqldb.user: root/; s/mysqldb.password:.*/mysqldb.password: dtstack/' /easy-agent-server/example-config.yml &&
RUN ln -sf /dev/stdout /tmp/easyagent-server.log && ln -sf /dev/stderr /tmp/easyagent-server-error.log && ln -sf /dev/stdout /tmp/api.log

WORKDIR /easy-agent-server
COPY bin /easy-agent-server/bin

#CMD ["./easy-agent-server","-c","example-config.yml","--debug"]
CMD ./bootstrap.sh
EXPOSE 8889
EXPOSE 8890
