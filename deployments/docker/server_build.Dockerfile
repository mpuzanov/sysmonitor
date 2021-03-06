FROM golang:1.13-alpine as builder
ENV APP_NAME sysmonitor
WORKDIR /opt/${APP_NAME}
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/${APP_NAME} ./cmd/sysmonitor

FROM ubuntu:18.04
ENV APP_NAME sysmonitor
LABEL name=${APP_NAME} maintainer="Mikhail Puzanov <mpuzanov@mail.ru>" version="1"
WORKDIR /opt/${APP_NAME}
COPY --from=builder /opt/${APP_NAME}/bin/${APP_NAME} ./bin/
COPY --from=builder /opt/${APP_NAME}/configs/prod.yaml ./configs/

RUN apt-get update \
    && apt-get -y install sysstat \
    && apt-get -y install tzdata \
    && dpkg-reconfigure --frontend noninteractive tzdata \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 50051
ENTRYPOINT ["./bin/sysmonitor"]
CMD ["grpc_server","-c", "./configs/prod.yaml"]