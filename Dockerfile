FROM golang:1.15-alpine AS STAGE_BUILD

ARG SERVICE

RUN apk update && apk add make git

WORKDIR /go/src/github.com/freemiumvpn/${SERVICE}/

ADD go.sum go.mod Makefile /go/src/github.com/freemiumvpn/${SERVICE}/
RUN make install

ADD . /go/src/github.com/freemiumvpn/${SERVICE}/

RUN make test
RUN make build
RUN mv ./${SERVICE} /${SERVICE}

FROM alpine:3.10 AS STAGE_SERVE

RUN apk update

ARG SERVICE

WORKDIR /work_dir


COPY --from=STAGE_BUILD /${SERVICE} /work_dir/service

ENTRYPOINT /work_dir/service
