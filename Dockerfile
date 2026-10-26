FROM golang:1.27.3-alpine@sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c as builder

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


FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a

LABEL maintainer="Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/restis /usr/local/bin/restis

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/restis"]
CMD ["--"]

# vim: set ft=Dockerfile:
