FROM golang:1.14

ENV NAME "integration_tests"

WORKDIR /opt/${NAME}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
CMD ["go", "test", "./tests/integration"]