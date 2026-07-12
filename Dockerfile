FROM golang:1.26.5-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 as builder

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


FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

LABEL maintainer="Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/restis /usr/local/bin/restis

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/restis"]
CMD ["--"]

# vim: set ft=Dockerfile:
