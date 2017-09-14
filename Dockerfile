FROM alpine:3.5

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
	&& apk add --update --no-cache jq groff less python \
	&& apk add --virtual .build-deps py-pip \
	&& pip install awscli \
	&& apk del .build-deps

COPY ./redis_auth_proxy /usr/bin/redis_auth_proxy
COPY ./docker-entrypoint.sh /docker-entrypoint.sh


ENTRYPOINT ["/docker-entrypoint.sh"]
