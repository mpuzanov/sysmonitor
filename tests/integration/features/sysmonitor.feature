# Тестирование сервиса sysmonitor
# тестирую факт потока статистики

Feature: service sysmonitor
    In order to test the sysmonitor application
    As an GRPC client operates with the service through API
	The service should be able to do the following

    Scenario: should test method Sysinfo
        When grpc client call method Sysinfo
        Then The error should be nil
        And The response data stream is not empty