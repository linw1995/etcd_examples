# ETCD examples

Code examples about interaction with ETCD.

## Usage

```shell
# Launch a local ETCD server
podman run \
	--name etcd \
	--publish 2379:2379 \
	--env ALLOW_NONE_AUTHENTICATION=yes \
	--env ETCD_ADVERTISE_CLIENT_URLS=http://localhost:2379 \
	-d \
	bitnami/etcd:3.4.15

# Shutdown
podman rm -f etcd

# Purge all key-values
podman run --rm --net=host bitnami/etcd etcdctl del / --prefix
```

### [lock_get_put/main.go](lock_get_put/main.go) lock, read, then modify the value of key

```shell
# Watch all key-values revisions
podman run --rm --net=host bitnami/etcd etcdctl watch /bar --prefix

# Run
ETCD_CLUSTER=http://localhost:2379 go run ./lock_get_put 
```

### [lock_get_put_v2/collections/value_test.go](lock_get_put_v2/collections/value_test.go) modify ETCD value with lock

```shell
# Watch all key-values revisions
podman run --rm --net=host bitnami/etcd etcdctl watch / --prefix

# Run
ETCD_CLUSTER=http://localhost:2379 go test ./lock_get_put_v2/... -v
```

### [value_history](value_history/main.go) Get all historical modifications of value

```shell
# Run
ETCD_CLUSTER=http://localhost:2379 go run ./value_history
```
