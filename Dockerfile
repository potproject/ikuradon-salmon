FROM golang:1.13.12-alpine AS build-env

ENV GO111MODULE=on

RUN apk --no-cache add git make build-base

WORKDIR /go/src/app

COPY . .

RUN mkdir -p /build
RUN go build -a -tags "netgo" -installsuffix netgo -ldflags="-s -w -extldflags \"-static\"" -o=/build/app main.go

FROM alpine:latest
RUN apk --no-cache add tzdata
#RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

COPY --from=build-env /build/app /build/app

RUN chmod u+x /build/app && mkdir /data

ENTRYPOINT ["/build/app"]