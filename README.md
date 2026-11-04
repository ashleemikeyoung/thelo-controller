# Thelo Controller

The Thelo Controller (**`theloctl`**) is the control server of Qavats Ranch, LLC's private
networks. It runs a **thelonet**: it signs devices in, hands out their addresses and keys,
and applies the network's policy. Devices connect with [Thelo](https://github.com/ashleemikeyoung/tailscale),
our client, or with the stock Tailscale apps.

This repository is a fork of [Headscale](https://github.com/juanfont/headscale). It is not
affiliated with or endorsed by the Headscale project or Tailscale Inc.; "Tailscale" is Tailscale
Inc.'s trademark. Headscale's own README is kept as [README.upstream.md](README.upstream.md).

## Branches

| Branch | What it is |
|---|---|
| `thelo` (default) | Headscale's development line with our commits on top |
| `main` | a plain mirror of upstream `main`, used to pick up fixes |

## What is different from Headscale

- **Access that follows the identity provider:** sign-in through Keycloak (OIDC), and a periodic
  check against Keycloak's admin API. People who are deleted, disabled or lose the network's role
  have their devices expired within minutes (`oidc.keycloak` in the config).
- **The names people read:** the command is `theloctl`; help, messages and the web pages say Thelo
  and the Thelo Controller, and a tailnet is a thelonet. Real paths and identifiers in error messages
  are left as they are. The brand layer is `cmd/headscale/cli/brand.go`.
- **Default paths** under `thelo-controller` (`/etc/thelo-controller`, `/var/lib/thelo-controller`).
- **No release check against GitHub:** updates come from our own builds.
- **Unchanged, on purpose:** the protocol, the database schema, the API and config keys, and the
  `HEADSCALE_*` environment variables, so an existing Headscale moves over with its nodes intact.

## Building

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w" -o theloctl ./cmd/headscale
```

## Ownership and licence

The Thelo Controller is the work and property of Qavats Ranch, LLC. Copyright © 2026 Qavats Ranch, LLC.

It is free software under the **GNU Affero General Public License, version 3 or later**
([LICENSE](LICENSE)); its web pages link to the source of the version they run, as the AGPL asks.
Code that comes from Headscale keeps its BSD 3-Clause licence
([LICENSE-BSD-3-Clause](LICENSE-BSD-3-Clause)); [NOTICE](NOTICE) explains which covers what. The
names and marks Thelo, Qavats and ElRoi, including the logos in `hscontrol/assets/brand`, belong to
Qavats Ranch, LLC and are not covered by the licence.
