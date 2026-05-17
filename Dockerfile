FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETARCH
ARG CADDY_VERSION=latest
ARG BUNDLE=latest

RUN apk add --no-cache jq && \
    go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest

COPY modules.json /modules.json

RUN set -eux; \
    MODULES=$(jq -r --arg b "$BUNDLE" \
      '.bundles[] | select(.name == $b) | .modules[] | "--with \(.import)@\(.version)"' \
      /modules.json | tr '\n' ' '); \
    if [ "$CADDY_VERSION" = "latest" ]; then \
      CGO_ENABLED=0 GOARCH="$TARGETARCH" xcaddy build $MODULES --output /usr/local/bin/caddy; \
    else \
      CGO_ENABLED=0 GOARCH="$TARGETARCH" xcaddy build "$CADDY_VERSION" $MODULES --output /usr/local/bin/caddy; \
    fi

FROM caddy:2-alpine
COPY --from=builder /usr/local/bin/caddy /usr/bin/caddy
