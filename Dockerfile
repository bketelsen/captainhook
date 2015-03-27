FROM golang:1.4.2-onbuild
VOLUME /config
CMD app -configdir /config
