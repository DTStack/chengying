FROM nginx:1.13.8-alpine

ENV \
  OSS_EXTERNAL=dtstack-download.oss-cn-hangzhou.aliyuncs.com \
  OSS_INTERNAL=dtstack-download.oss-cn-hangzhou-internal.aliyuncs.com
  
################# update apk repositories #################
RUN \
  echo "" > /etc/apk/repositories && \
  echo "http://mirrors.aliyun.com/alpine/v3.8/main" >> /etc/apk/repositories && \
  echo "http://mirrors.aliyun.com/alpine/v3.8/community" >> /etc/apk/repositories && \
  apk update
################# update apk repositories #################

RUN apk --update add curl netcat-openbsd

RUN cd /tmp \
  && curl -L http://${OSS_EXTERNAL}/insight/em/2.5/managefront/managefront-2.5.tgz > managerfront.tgz \
  && tar -zxf managerfront.tgz && mv dist/ /usr/share/nginx/html/easymanager && rm -f managerfront.tgz

COPY health_check.sh /dtstack/
COPY nginx.default.conf /etc/nginx/conf.d/default.conf
COPY drpeco /dtstack/

CMD [ "./dtstack/entrypoint.sh" ]
