#!/bin/sh

export DEMO_PRIVATEKEY=/home/ubuntu/.lego/certificates/lab.owulveryck.info.key
export DEMO_CERTIFICATE=/home/ubuntu/.lego/certificates/lab.owulveryck.info.crt
export DEMO_PORT=8080
go build demo.go
sudo setcap CAP_NET_BIND_SERVICE=+eip /home/ubuntu/GOPROJECTS/src/github.com/owulveryck/
./demo
