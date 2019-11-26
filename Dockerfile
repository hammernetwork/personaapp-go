FROM golang:1.12

RUN mkdir -p /go/bitbucket/personaapp

WORKDIR /go/bitbucket.org/personaapp

COPY . /go/bitbucket.org/personaapp

RUN make build