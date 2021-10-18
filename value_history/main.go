package main

import (
	"context"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	etcdError "github.com/coreos/etcd/error"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/linw1995/etcd_examples"
)

const (
	VALUE_KEY = "/char"
)

func GetKeyHistory(ctx context.Context, client *clientv3.Client, key string) (rv []*mvccpb.KeyValue, err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	kvc := clientv3.NewKV(client)
	resp, err := kvc.Get(ctx, key)
	if err != nil {
		return
	}
	if resp.Count == 0 {
		err = etcdError.NewRequestError(etcdError.EcodeKeyNotFound, key)
		return
	}
	rev := resp.Kvs[0].CreateRevision

	watcher := clientv3.NewWatcher(client)
	select {
	case resp := <-watcher.Watch(ctx, key, clientv3.WithRev(rev)):
		for _, e := range resp.Events {
			rv = append(rv, e.Kv)
		}
		return
	case <-ctx.Done():
		err = ctx.Err()
		return
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cli, err := clientv3.New(etcd_examples.EtcdCfgFromEnv())
	if err != nil {
		log.Printf("clientv3.New err: %s\n", err)
		return
	}
	defer cli.Close()

	kvc := clientv3.NewKV(cli)
	for _, v := range []string{"a", "b", "c", "d", "e"} {
		_, err = kvc.Put(ctx, VALUE_KEY, v)
		if err != nil {
			log.Printf("kv.Put err: %s\n", err)
			return
		}
	}

	kvs, err := GetKeyHistory(ctx, cli, VALUE_KEY)
	if err != nil {
		log.Printf("GetKeyHistory err: %s\n", err)
		return
	}
	for idx, kv := range kvs {
		log.Printf("Key#%d: %s\n", idx, kv)
	}

	_, err = GetKeyHistory(ctx, cli, "unknown")
	if err != nil {
		log.Printf("GetKeyHistory err: %s\n", err)
		return
	}
}
