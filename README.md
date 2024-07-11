# v-agent

## Overview
`v-agent` is a daemon that gathers metrics and writes them to a compatible [remote_write](https://prometheus.io/docs/specs/remote_write_spec/) endpoint.

It acts like prometheus (scraping metrics), re-labeling metrics, then pushing them off to a remote_write compatible endpoint. Internally, we use mimir and write our metrics there but other implementations _should_ work.

### Metrics
All metrics that are specifically created with `v-agent` are prefixed with `v_`. Scraped metrics are not modified other than the addition of labels.

Every metric will have all metrics in `labels_config` added to it. The following are special labels:
- `hostname`: Pulled automatically. Set with `HOSTNAME` environment variable or `os.Hostname()`
- `subid`: The subscription ID for the underlying service. Can be set in `config.yaml`. If it's not set an attempt is made to pull it from metadata API.
- `vpsid`: The VPS ID. Can be set in `config.yaml`. If it's not set an attempt is made to pull it from the metadata API.

The above labels are added to ensure that the metric is unique. A lack of uniqueness can result in metrics getting overwritten/clobbered.

System metrics that are collected:
- CPU utilization: system, user, steal, utilization, etc.
- Memory utilization: cached, buffered, utilization, etc.
- Load average: 1, 5, 15, and tasks
- Disk stats: writes/reads, etc.
- Filesystem stats: bytes, inodes, utilization
- NIC: bytes, packets, errors, etc.

Kubernetes:
- `v_kube_apiserver_healthy` that is `0` (if healthy) or `1` if not healthy based on response from kube-apiserver `/healthz` endpoint.
- Every metric from `/metrics`

Etcd:
- `v_etcd_healthy` that is `0` (if healthy) or `1` if not healthy based on response from etcd `/health` endpoint.
- Every metric from `/metrics`

Konnectivty:
- `v_konnectivity_healthy` that is `0` (if healthy) or `1` if not healthy based on response from konnectivity `/healthz` endpoint.
- Every metric from `/metrics`

HAProxy:
- `v_haproxy_healthy` that is `0` (if healthy) or `1` if not healthy based on response from `/metrics` endpoint.
- Every metric from `/metrics`

Ceph:
- `v_ceph_healthy`: Not implemented yet.
- Every metric from `/metrics`

## Usage
Configuration is through `config.yaml`, sample:

```yaml
debug: true                      # debug output
interval: 60                     # interval to scrape metrics
endpoint: https://endpoint...    # remote endpoint
basic_auth_user: ""              # basic auth user
basic_auth_pass: ""              # basic auth pass
check_vendor: false              # when true, vendor must be "Vultr"; set to false otherwise
labels_config:                   # any labels below will be added to all metrics
  hostname: ""                   # empty string uses local hostname, unset (nil) doesnt use, non-empty string uses specified label
  subid: ""                      # empty string pulls from userdata, unset (nil) doesnt use, non-empty string uses specified label
  vpsid: ""                      # empty string pulls from userdata, unset (nil) doesnt use, non-empty string uses specified label
  product: vke                   # unset (nil) doesnt use, non-empty string uses specified label. Note: This label is used to determine subid for vke/vlb/vfs
  any: any                       # any key/value label
probes_api:
  listen: 0.0.0.0
  port: 7091
metrics_config:
  agent:
    load_avg:
      enabled: true
    cpu:
      enabled: true
    memory:
      enabled: true
    nic:
      enabled: true
    disk_stats:
      enabled: true
      filter: "sr0" # regex
    file_system:
      enabled: true
    kubernetes:
      enabled: true
      endpoint: https://localhost:6443
      kubeconfig: /var/lib/kubernetes/admin.kubeconfig
    konnectivity:
      enabled: true
      metrics_endpoint: http://localhost:8133 # /metrics
      health_endpoint: http://localhost:8092 # /healthz
    etcd:
      enabled: true
      cacert: /var/lib/kubernetes/ca.pem
      cert: /var/lib/kubernetes/kubernetes.pem
      key: /var/lib/kubernetes/kubernetes-key.pem
      endpoint: https://10.1.96.3:2379 # /metrics
    nginx_vts:
      enabled: false
      endpoint: http://localhost:9001 # /metrics
    v_cdn_agent:
      enabled: false
      endpoint: http://localhost:9093 # /metrics
    haproxy:
      enabled: false
      endpoint: http://localhost:8404 # /metrics
    ceph:
      enabled: false
      endpoint: http://localhost:9283 # /metrics
    v_dns:
      enabled: false
      endpoint: http://localhost:9053 # /metrics
    smart:
      enabled: false
      block_devices: # must exist, if not set, block devices are used from /sys/block/ (except for dmX and loopX)
      - /dev/sda
  kubernetes:
    pods: # v-agent must be running inside k8s for this to work
      enabled: false
      namespaces:
      - rook-ceph
      - default
    dcgm: # v-agent must be running inside k8s for this to work
      enabled: false
      namespace: gpu-operator        # namespace
      endpoint: nvidia-dcgm-exporter # name of the endpoint: k get endpoints
```

Currently, it's not 100% compatible with Kubernetes, that is to say running inside k8s and able to scrape k8s metrics. Right now it's largely an agent used to scrape metrics for services.

## Building
Note: Agent must be built with cgo disabled, not doing so will result in GLIBC errors being thrown: `CGO_ENABLED=0 go build -o v-agent cmd/v-agent/main.go`