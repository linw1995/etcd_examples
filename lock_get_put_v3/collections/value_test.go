package collections

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"

	"github.com/linw1995/etcd_examples"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	client *clientv3.Client
)

func init() {
	var err error
	client, err = clientv3.New(etcd_examples.EtcdCfgFromEnv())
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err = clientv3.NewKV(client).Delete(ctx, "/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
}

func TestCommon(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	v := ETCDValue{
		Key:    "/bar",
		Client: client,
	}

	t.Run("value not exists", func(t *testing.T) {
		raw, err := v.Get(ctx)
		if raw != nil {
			t.Errorf("get cannot return: %s", raw)
		}
		if err != ErrNotFound {
			t.Fatalf("wrong err: %s", err)
		}
	})

	payload := ([]byte)("abc")

	t.Run("put value", func(t *testing.T) {
		err := v.Put(ctx, payload)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get value", func(t *testing.T) {
		raw, err := v.Get(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if raw != nil && !bytes.Equal(raw, payload) {
			t.Errorf("get cannot return: %s", raw)
		}
	})
}

func TestCommonWithDelete(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	v := ETCDValue{
		Key:    "/aaa",
		Client: client,
	}
	payload := ([]byte)("abc")

	t.Run("put value", func(t *testing.T) {
		err := v.Put(ctx, payload)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("get value", func(t *testing.T) {
		raw, err := v.Get(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if raw != nil && !bytes.Equal(raw, payload) {
			t.Errorf("get cannot return: %s", raw)
		}
	})

	t.Run("delete value", func(t *testing.T) {
		err := v.Del(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("value not exists", func(t *testing.T) {
		raw, err := v.Get(ctx)
		if raw != nil {
			t.Errorf("get cannot return: %s", raw)
		}
		if err != ErrNotFound {
			t.Fatalf("wrong err: %s", err)
		}
	})
}

func TestWithLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	wg := sync.WaitGroup{}

	var errs []error
	for i := 0; i < 2; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			v := ETCDValue{
				Key:    "/boo",
				Client: client,
			}

			err := v.WithLock(ctx, func(c context.Context) error {
				t.Log("running", time.Now())
				time.Sleep(time.Second)
				raw, err := v.Get(c)

				var payload string
				if err == ErrNotFound {
					payload = "1"
				} else if err == nil {
					payload = string(raw) + "-2"
				} else {
					return err
				}
				err = v.Put(c, []byte(payload))
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				errs = append(errs, err)
			}
		}()
	}
	wg.Wait()

	if len(errs) > 0 {
		t.Fatal(errs)
	}
}
