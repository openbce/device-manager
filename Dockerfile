FROM golang:1.18.3 as builder

ENV GOPATH /opt/

ADD . /opt/src/openbce.io/device-manager

RUN apt-get update && apt-get -y install libibverbs-dev build-essential
RUN cd /opt/src/openbce.io/device-manager/ && make


FROM ubuntu:18.04

RUN apt-get update && apt-get -y install libibverbs1
COPY  --from=builder /opt/src/openbce.io/device-manager/_output/device-manager /opt/bce-device-manager

ENTRYPOINT ["/opt/bce-device-manager"]
