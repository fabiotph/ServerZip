FROM golang:1.17


RUN mkdir -p /server

COPY . /server

WORKDIR /server
RUN go build -o docker-test
WORKDIR /server
RUN ./docker-test
