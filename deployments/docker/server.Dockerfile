FROM ubuntu:18.04
ENV APP_NAME sysmonitor
LABEL name=${APP_NAME} maintainer="Mikhail Puzanov <mpuzanov@mail.ru>" version="1"
WORKDIR /opt/${APP_NAME}

COPY ./sysmonitor ./bin/
COPY ./configs/prod.yaml ./configs/

RUN apt-get update \
    && apt-get -y install sysstat \
    && apt-get -y install tzdata \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 50051
ENTRYPOINT ["./bin/sysmonitor"]
CMD ["grpc_server","-c", "./configs/prod.yaml"]