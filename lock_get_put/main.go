package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/linw1995/etcd_examples"
)

const (
	// The key of value can not be a prefix of other keys.
	// It will cause lock forever.
	VALUE_KEY = "/bar"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	cli, err := clientv3.New(etcd_examples.EtcdCfgFromEnv())
	if err != nil {
		log.Printf("clientv3.New err: %s\n", err)
		return
	}
	defer cli.Close()

	task := func(idx int) {
		session, err := concurrency.NewSession(cli, concurrency.WithContext(ctx))
		if err != nil {
			log.Printf("concurrency.NewSession err: %s\n", err)
			return
		}

		// Lock
		locker := concurrency.NewMutex(session, VALUE_KEY)
		err = locker.Lock(ctx)
		if err != nil {
			log.Printf("locker.Lock err: %s\n", err)
			return
		}
		defer func() {
			err := locker.Unlock(context.Background())
			if err != nil {
				log.Printf("locker.Unlock err: %s\n", err)
			}
		}()

		kv := clientv3.NewKV(cli)

		// Read
		resp, err := kv.Get(ctx, VALUE_KEY)
		if err != nil {
			log.Printf("kv.Get err: %s", err)
			return
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
			log.Printf("kv.Put err: %s\n", err)
			return
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
		log.Printf("kv.Get err: %s\n", err)
		return
	}
	fmt.Println((string)(resp.Kvs[0].Value))
}
