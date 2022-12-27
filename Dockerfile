FROM golang:alpine as builder

COPY . /src/restis
WORKDIR /src/restis

ENV CGO_ENABLED=0

RUN set -ex \
 && apk add --update \
      build-base \
      git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath

FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/restis /usr/local/bin/restis

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/restis"]
CMD ["--"]

# vim: set ft=Dockerfile:
