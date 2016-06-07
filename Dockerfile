FROM golang:1.6.2-alpine

RUN apk add --update git && apk add --update make && rm -rf /var/cache/apk/*

ADD . /go/src/github.com/r3labs/workflow-manager
WORKDIR /go/src/github.com/r3labs/workflow-manager

RUN make deps && go install

ENTRYPOINT /go/bin/workflow-manager
