FROM --platform=$BUILDPLATFORM golang:alpine AS builder
ARG TARGETARCH
ARG CADDY_VERSION=latest
ARG BUNDLE=latest

RUN go install github.com/caddyserver/xcaddy/cmd/xcaddy@v0.4.6

COPY modules.json /modules.json
COPY go.mod /src/go.mod
COPY tools/modtool /src/tools/modtool

RUN set -eux; \
    MODULES=$(cd /src && go run ./tools/modtool -manifest /modules.json xcaddy-args "$BUNDLE"); \
    if [ "$CADDY_VERSION" = "latest" ]; then \
      CGO_ENABLED=0 GOARCH="$TARGETARCH" xcaddy build $MODULES --output /usr/local/bin/caddy; \
    else \
      CGO_ENABLED=0 GOARCH="$TARGETARCH" xcaddy build "$CADDY_VERSION" $MODULES --output /usr/local/bin/caddy; \
    fi

FROM caddy:2-alpine
COPY --from=builder /usr/local/bin/caddy /usr/bin/caddy
