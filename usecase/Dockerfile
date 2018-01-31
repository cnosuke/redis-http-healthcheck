FROM golang:alpine AS builder
RUN apk --no-cache -u add git make curl && \
  go get github.com/cnosuke/redis-http-healthcheck && \
  cd $GOPATH/src/github.com/cnosuke/redis-http-healthcheck && \
  make deps && make cross-build && \
  cp ./dist/linux-amd64/healthchecker /.

FROM redis:4-alpine
COPY --from=builder /healthchecker /usr/bin/.

ENV HEALTHCHECKER_CONFIG_PATH /tmp/config.yml
ADD docker-entrypoint.sh /usr/local/bin/
