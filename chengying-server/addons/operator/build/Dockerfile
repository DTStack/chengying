FROM debian:buster-slim

ENV OPERATOR=/usr/local/bin/mole-operator \
    USER_UID=1001 \
    USER_NAME=mole-operator

# install operator binary
COPY build/_output/bin/mole-operator ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
