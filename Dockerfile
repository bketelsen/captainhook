# Docker multi-stage build file
# Requires docker 17.05 or newer.

FROM golang:1.4.2-onbuild as builder
VOLUME /config
CMD app -listen-addr 0.0.0.0:8080 -configdir /config


FROM ubuntu:16.04
COPY --from=builder /go/bin/app /usr/local/bin/captainhook
ENTRYPOINT ["/usr/local/bin/captainhook"]
