# v-agent

## Overview
`v-agent` is a daemon that gathers metrics and writes them to a `v-proxy`.

`v-proxy` takes a `v-agent` remote write and authenticates it, if it's valid it sends it off to mimir.

### Architecture
![](./docs/metrics-arch.drawio.png)

`v-agent` must auth with `v-proxy`, this is done passively via headers:
- `X-Vultr-SubID`: Pulled from userdata.
- `X-Vultr-Key`: Pulled from userdata.

`v-proxy` will verify both the headers; once they're verified will write metrics to mimir.

Authentication is passive via headers; if the headers are missing or invalid then request is simply not forwarded.

## Tools
- `v-agent`: Gathers and sends metrics to `v-proxy`.
- `v-proxy`: Receives metrics from `v-agent` and sends to mimir.

## Configuration
Both have a `config.yaml` file. Both have CLI switches. Both configurations can be overridden with envionment variables.
