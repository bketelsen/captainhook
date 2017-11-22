FROM golang as builder
RUN go get -u github.com/golang/dep/cmd/dep
RUN go env
ADD . /go/src/github.com/bketelsen/captainhook
WORKDIR /go/src/github.com/bketelsen/captainhook
RUN make clean
RUN make
VOLUME /config
CMD /go/src/github.com/bketelsen/captainhook/bin/captainhook


FROM ubuntu:16.04
COPY --from=builder /go/bin/app /usr/local/bin/captainhook
ENTRYPOINT ["/usr/local/bin/captainhook"]
