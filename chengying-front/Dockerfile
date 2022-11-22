FROM nginx:alpine

LABEL maintainer="guyan@dtstack.com"

COPY dist /usr/share/nginx/html/easymanager
COPY dt-alert.conf /usr/share/nginx/html/dt-alert
COPY em.conf dt-alert.conf /etc/nginx/conf.d/

RUN rm /etc/nginx/conf.d/default.conf

EXPOSE 80
EXPOSE 8600