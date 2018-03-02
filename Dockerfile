FROM golang:1.10
ADD . /go/src/github.com/bketelsen/captainhook
WORKDIR /go/src/github.com/bketelsen/captainhook
RUN go get -u golang.org/x/vgo
RUN CGO_ENABLED=0 GOOS=linux vgo build -a -installsuffix cgo -o captainhook .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/bketelsen/captainhook/captainhook .
VOLUME /config
CMD ["./captainhook"]



