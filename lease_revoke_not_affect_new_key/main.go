package main

import (
	"context"
	"log"
	"time"

	"github.com/linw1995/etcd_examples"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	cli, err := clientv3.New(etcd_examples.EtcdCfgFromEnv())
	if err != nil {
		log.Printf("clientv3.New err: %s\n", err)
		return
	}
	defer cli.Close()

	session, err := concurrency.NewSession(cli, concurrency.WithContext(ctx))
	if err != nil {
		log.Printf("NewSession err: %s\n", err)
		return
	}

	putRes, err := cli.Put(ctx, "/abc", "with_lease", clientv3.WithLease(session.Lease()), clientv3.WithPrevKV())
	if err != nil {
		log.Printf("Put err: %s\n", err)
		return
	}
	log.Println(putRes)

	putRes, err = cli.Put(ctx, "/abc", "without_lease", clientv3.WithPrevKV())
	if err != nil {
		log.Printf("Put err: %s\n", err)
		return
	}
	log.Println(putRes)

	err = session.Close()
	if err != nil {
		log.Printf("session.Close err: %s\n", err)
		return
	}

	<-time.After(time.Second * 66)

	getRes, err := cli.Get(ctx, "/abc")
	if err != nil {
		log.Printf("Get err: %s\n", err)
		return
	}
	log.Println(getRes)
}
