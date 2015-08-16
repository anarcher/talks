FROM golang:1.4.2

RUN go get golang.org/x/tools/cmd/present

WORKDIR /go/src/slide

ADD . /go/src/slide
ENV GOPATH /go/src/slide/gokit/Godeps/_workspace:$GOPATH

EXPOSE 3999

ENTRYPOINT ["present"]
CMD ["-orighost=0.0.0.0"]

