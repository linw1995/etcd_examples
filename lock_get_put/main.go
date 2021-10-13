package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"

	"github.com/linw1995/etcd_examples"
)

const (
	VALUE_KEY = "/bar"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cli, err := clientv3.New(etcd_examples.EtcdCfgFromEnv())
	if err != nil {
		panic(err)
	}

	task := func(idx int) {
		session, err := concurrency.NewSession(cli)
		if err != nil {
			panic(err)
		}

		// Lock
		locker := concurrency.NewLocker(session, VALUE_KEY)
		locker.Lock()
		defer locker.Unlock()

		kv := clientv3.NewKV(cli)

		// Read
		resp, err := kv.Get(ctx, VALUE_KEY)
		if err != nil {
			panic(err)
		}

		// Generate new value
		var value string
		if len(resp.Kvs) == 0 {
			value = fmt.Sprint(idx)
		} else {
			value = string(resp.Kvs[0].Value) + fmt.Sprintf("-%d", idx)
		}

		// Write
		_, err = kv.Put(ctx, VALUE_KEY, value)
		if err != nil {
			panic(err)
		}
	}

	// Testing
	wg := sync.WaitGroup{}
	run := func(idx int) {
		wg.Add(1)
		defer wg.Done()
		task(idx)
	}

	for idx := 0; idx < 100; idx++ {
		go run(idx)
	}
	wg.Wait()

	kv := clientv3.NewKV(cli)
	resp, err := kv.Get(ctx, VALUE_KEY)
	if err != nil {
		panic(err)
	}
	fmt.Println((string)(resp.Kvs[0].Value))
}
