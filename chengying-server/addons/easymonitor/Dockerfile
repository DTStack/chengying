FROM lgatica/openssh-client

LABEL maintainer="hx@dtstack.com"

COPY easymonitor monitor.yml entrypoint.sh /easymonitor/

WORKDIR /easymonitor

CMD ./entrypoint.sh
