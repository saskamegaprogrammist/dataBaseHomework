FROM golang:1.13 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN ls
RUN go build .

FROM ubuntu:18.04 AS release

USER root

EXPOSE 5000

COPY --from=build /opt/app/dataBaseHomework /usr/bin/

CMD dataBaseHomework

