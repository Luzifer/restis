FROM golang:1.27.4-alpine@sha256:0a5f7fe796599ff487605987432e7619ae45beeb2fba4c2ce0b907879932aebf as builder

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


FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf

LABEL maintainer="Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/restis /usr/local/bin/restis

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/restis"]
CMD ["--"]

# vim: set ft=Dockerfile:
