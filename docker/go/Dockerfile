FROM golang:1.21 as base

FROM base as dev

RUN apt update && apt install -y \
     traceroute \
     net-tools \
     iproute2 \
     iptables \
     tcpdump \
     iputils-ping \
     telnet

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

WORKDIR /opt/app
CMD ["air"]