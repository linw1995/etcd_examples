package etcd_examples

import (
	"fmt"
	"os"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	ENV_ETCD_ENDPOINT = "ETCD_CLUSTER"
	ENV_ETCD_PASSWORD = "ETCD_PASSWORD"
	ENV_ETCD_USERNAME = "ETCD_USERNAME"
)

func EtcdCfgFromEnv() clientv3.Config {
	endpointRaw := os.Getenv(ENV_ETCD_ENDPOINT)
	if endpointRaw == "" {
		panic(fmt.Errorf("Missing env var: %s", ENV_ETCD_ENDPOINT))
	}
	endpoints := strings.Split(endpointRaw, ",")
	return clientv3.Config{
		Endpoints: endpoints,
		Username:  os.Getenv(ENV_ETCD_USERNAME),
		Password:  os.Getenv(ENV_ETCD_PASSWORD),
	}
}
