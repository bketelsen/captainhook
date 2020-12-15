FROM golang:1.15
ADD . /go/src/github.com/bketelsen/captainhook
WORKDIR /go/src/github.com/bketelsen/captainhook
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o captainhook .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/bketelsen/captainhook/captainhook .
VOLUME /config
CMD ["./captainhook"]



