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

### `v-agent`
`v-agent` is a prometheus compatible remote write client.

All metrics that are specifically created with `v-agent` are prefixed with `v_`. Scraped metrics are not modified other than the addition of labels.

System metrics that are collected:
- CPU utilization: system, user, steal, utilization, etc.
- Memory utilization: cached, buffered, utilization, etc.
- Load average: 1, 5, 15, and tasks
- Disk stats: writes/reads, etc.
- Filesystem stats: bytes, inodes, utilization
- NIC: bytes, packets, errors, etc.

Kubernetes:
- `v_kube_apiserver_healthy` that is `1` (if healthy) or `0` if not healthy based on response from kube-apiserver `/healthz` endpoint.
- Every metric from `/metrics`

Etcd:
- `v_etcd_healthy` that is `1` (if healthy) or `0` if not healthy based on response from etcd `/health` endpoint.
- Every metric from `/metrics`

## Configuration
Both have a `config.yaml` file. Both have CLI switches. Both configurations can be overridden with envionment variables.

## Building
Note: Agent must be built with cgo disabled for VKE guests, not doing so will result in GLIBC errors being thrown: `CGO_ENABLED=0 go build -o v-agent cmd/v-agent/main.go`
