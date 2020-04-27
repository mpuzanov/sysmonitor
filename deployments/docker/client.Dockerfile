FROM alpine:3.11
ENV APP_NAME sysmonitor
LABEL name=${APP_NAME} maintainer="Mikhail Puzanov <mpuzanov@mail.ru>" version="1"
WORKDIR /opt/${APP_NAME}

COPY ./sysmonitor ./bin/

RUN apk add --no-cache tzdata \
    && apk add -U --no-cache ca-certificates \
    && adduser -D -g appuser appuser \
    && chmod -R 755 ./

USER appuser

ENTRYPOINT ["./bin/sysmonitor"]
CMD ["grpc_client"]
#CMD ["grpc_client","--server='0.0.0.0:50051'"]