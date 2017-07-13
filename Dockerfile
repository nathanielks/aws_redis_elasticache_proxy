FROM scratch
COPY ./redis_auth_proxy /redis_auth_proxy
ENTRYPOINT ["/redis_auth_proxy"]
