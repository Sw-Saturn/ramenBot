FROM golang:1.13.3-alpine3.10 AS build-bin

ENV GO111MODULE=on
RUN apk --no-cache add git make build-base
WORKDIR /go/src/github.com/Sw-Saturn/ramenBot
COPY . .

RUN mkdir -p /build
RUN go build -a -tags "netgo" -installsuffix netgo -ldflags="-s -w -extldflags \"-static\"" -o=/build/app main.go

FROM alpine:3.10.2

COPY .env .
RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
COPY --from=build-bin /build/app /build/app
RUN chmod u+x /build/app

EXPOSE 8080
ENTRYPOINT ["/build/app"]