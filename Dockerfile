FROM golang:1.15-alpine AS STAGE_BUILD

ARG SERVICE

RUN apk update && apk add make git

WORKDIR /go/src/github.com/freemiumvpn/${SERVICE}

ADD go.sum go.mod Makefile /go/src/github.com/freemiumvpn/${SERVICE}
RUN make install

ADD . /go/src/github.com/freemiumvpn/${SERVICE}

RUN make test
RUN make build

FROM alpine:3.10 AS STAGE_SERVE

ARG SERVICE

RUN mkdir /app

COPY --from=STAGE_BUILD /go/src/github.com/freemiumvpn/${SERVICE} /app/$SERVICE

ENTRYPOINT /app/$SERVICE
