FROM golang:alpine

RUN apk update
RUN apk add git

ARG CACHEBUST=1
COPY . /github/workspace
COPY .github/actions/gomodtidy/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
