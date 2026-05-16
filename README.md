# sharan-caddy

Custom Caddy images for [Sharan](https://github.com/bemoty/sharan)

## Bundles

| Tag | Contents |
|-----|----------|
| `:stock` | Upstream `caddy:2-alpine` re-tag. No extra modules. |
| `:cloud` | layer4, DNS-01 (Cloudflare, Route53, Hetzner, DigitalOcean). |
| `:homelab` | Same as `cloud`, plus Tailscale. |

Images are published to `ghcr.io/bemoty/sharan-caddy`. See the [Sharan docs](https://sharan.io) for details.
