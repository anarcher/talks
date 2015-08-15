FROM golang:1.4.2

RUN go get golang.org/x/tools/cmd/present

WORKDIR /go/src/gokit

ADD . /go/src/gokit

RUN cd gokit ; go get -v

EXPOSE 3999

CMD ["present","-orighost=localhost"]



