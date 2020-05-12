[![Build Status](https://travis-ci.org/mpuzanov/sysmonitor.svg?branch=master)](https://travis-ci.org/mpuzanov/sysmonitor)[![Go Report Card](https://goreportcard.com/badge/github.com/mpuzanov/sysmonitor)](https://goreportcard.com/report/github.com/mpuzanov/sysmonitor)
# Проектная работа "Системный мониторинг"

## Общее описание

Демон - программа, собирающая информацию о системе, на которой запущена, и отправляющая её своим клиентам по GRPC.

## Запуск сервера

Варианты запуска:
 - make run
 - ./sysmonitor grpc_server --config=configs/prod.yaml --port=50051
 - ./sysmonitor grpc_server --port=50051

Через файл можно указать, какие из подсистем сбора включены/выключены

## Запуск клиента

Варианты запуска:
 - make run-client
 - ./sysmonitor grpc_client --server=":50051"
 - ./sysmonitor grpc_client --server=":50051" --timeout=5 --period=15

## Тестирование сервиса

запуск:

    make test

поднятие сервера и клиента в докер контейнере

    make up

запуск интеграционных тестов

    make integration-tests

## Виды собираемой статистики на сервере

- load average
- загрузка CPU
- загрузка дисков
- top talkers по сети
