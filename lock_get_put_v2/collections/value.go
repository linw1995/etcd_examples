package collections

import (
	"context"
	"fmt"
	"log"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

type ETCDValue struct {
	Key    string
	Client *clientv3.Client
}

func (v *ETCDValue) Tx(ctx context.Context, next func(context.Context) error) error {
	session, err := concurrency.NewSession(v.Client, concurrency.WithContext(ctx))
	if err != nil {
		return err
	}
	lock := concurrency.NewMutex(session, v.Key)
	err = lock.Lock(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err = lock.Unlock(context.Background())
		if err != nil {
			log.Printf("ETCDValue(%s) Tx.Unlock err: %s", v.Key, err)
		}
	}()

	return next(context.WithValue(ctx, etcd_value_tx{}, 0))
}

var (
	ErrNotFound = fmt.Errorf("not found")
)

func (v *ETCDValue) Get(ctx context.Context) (data []byte, err error) {
	next := func(ctx context.Context) error {
		kv := clientv3.NewKV(v.Client)
		resp, err := kv.Get(ctx, v.Key)
		if err != nil {
			return err
		}
		if len(resp.Kvs) == 0 {
			return ErrNotFound
		}
		data = resp.Kvs[0].Value
		return nil
	}
	if ctx.Value(etcd_value_tx{}) == nil {
		err = v.Tx(ctx, next)
	} else {
		err = next(ctx)
	}
	return
}

func (v *ETCDValue) Put(ctx context.Context, value []byte) (err error) {
	next := func(ctx context.Context) error {
		kv := clientv3.NewKV(v.Client)
		_, err := kv.Put(ctx, v.Key, string(value))
		return err
	}
	if ctx.Value(etcd_value_tx{}) == nil {
		err = v.Tx(ctx, next)
	} else {
		err = next(ctx)
	}
	return
}

func (v *ETCDValue) Del(ctx context.Context) (err error) {
	next := func(ctx context.Context) error {
		kv := clientv3.NewKV(v.Client)
		_, err := kv.Delete(ctx, v.Key)
		return err
	}
	if ctx.Value(etcd_value_tx{}) == nil {
		err = v.Tx(ctx, next)
	} else {
		err = next(ctx)
	}
	return
}

type etcd_value_tx struct{}
