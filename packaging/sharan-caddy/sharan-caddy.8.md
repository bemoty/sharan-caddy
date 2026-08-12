% "SHARAN-CADDY" "8" "" "sharan-caddy" "System Manager's Manual"

# NAME

sharan-caddy - Caddy web server built for use with Sharan

# SYNOPSIS

**sharan-caddy** *command* [*flags*]

# DESCRIPTION

sharan-caddy is a build of the Caddy web server. It is Caddy; the only
difference is the set of bundled modules, which is chosen to match what
Sharan relies on. It installs under its own binary name, unit name, and
configuration paths, so it does not conflict with a distribution **caddy**
package.

# MODULES

sharan-caddy bundles the following modules in addition to Caddy's standard
module set:

layer4, tailscale, dns.providers.cloudflare, dns.providers.route53,
dns.providers.hetzner, dns.providers.ionos, dns.providers.digitalocean

# FILES

*/etc/sharan-caddy/Caddyfile*
: Configuration. Ships containing only an admin endpoint block.

*/var/lib/sharan-caddy*
: Data directory, owned by the service user.

*/usr/bin/sharan-caddy*
: The binary.

# SERVICE

sharan-caddy runs as the **sharan-caddy.service** systemd unit, under the
sharan-caddy user and group. The unit installs disabled and must be enabled
explicitly. Privileged port binding comes from the unit's ambient
capability, not from file capabilities on the binary.

# SEE ALSO

Run **sharan-caddy help** for the full command reference. Caddy's own
documentation is at https://caddyserver.com/docs/.

A distribution package named **caddy** may also provide a **caddy(8)**
manual page; that is a different binary without these modules.
