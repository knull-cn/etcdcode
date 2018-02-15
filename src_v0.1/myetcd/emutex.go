package myetcd

import (
	"github.com/coreos/etcd/clientv3/concurrency"
	"golang.org/x/net/context"
	"path/filepath"
	"sync"
)

type EMutex struct {
	emtx *concurrency.Mutex
	mtx  sync.Mutex
}

func NewEMutex(eaddr []string, key string) *EMutex {
	cli, ok := NewClientObj(eaddr)
	if !ok {
		return nil
	}
	sess := newSession(cli)
	if sess == nil {
		cli.Close()
		return nil
	}
	emutex := concurrency.NewMutex(sess, filepath.Join(key, "mtx"))
	return &EMutex{
		emtx: emutex,
	}
}

func (em *EMutex) Lock() error {
	em.mtx.Lock()
	err := em.emtx.Lock(context.TODO())
	em.mtx.Unlock()
	return err
}

func (em *EMutex) Unlock() error {
	em.mtx.Lock()
	err := em.emtx.Unlock(context.TODO())
	em.mtx.Unlock()
	return err
}
