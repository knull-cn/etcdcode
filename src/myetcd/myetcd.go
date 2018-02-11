package myetcd

import (
	"fmt"
	clientv3 "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"logger"
	//"path/filepath"
	"time"
)

type ETCD_EVENT int32

const (
	EE_PUT ETCD_EVENT = 0
	EE_DEL            = 1
)

const (
	timeout = time.Second * 5
)

// type EtcdInfo struct {
// 	SysName  string   //系统名称;
// 	EtcdAddr []string //系统ETCD;
// 	//TLS:似乎不需要,因为都是弄内网,不给外网用;外网提供服务器接口;
// }

type MyEtcd struct {
	*clientv3.Client
	prefix string
	addres []string
	//cli    *clientv3.Client
}

type EtcdKey []byte
type EtcdValue []byte
type EtcdKeyValue struct {
	Key   EtcdKey
	Value EtcdValue
}

func (ekv *EtcdKeyValue) String() string {
	return fmt.Sprintf("key=%s;value=%s", ekv.Key, ekv.Value)
}

func NewMyEtcd(SysName string, EtcdAddr []string) *MyEtcd {
	return &MyEtcd{
		nil,
		"/" + SysName + "/",
		EtcdAddr,
	}
}

func (me *MyEtcd) Initialize() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   me.addres,
		DialTimeout: 1 * time.Second,
	})
	if err != nil {
		return err
	}
	me.Client = cli
	//clean;
	del(me.Client, me.prefix)
	//
	return nil
}

func (me *MyEtcd) KeepAlive(leaseid clientv3.LeaseID) error {
	kachan, err := me.Client.KeepAlive(context.TODO(), leaseid)
	go func() {
		for {
			rsp := <-kachan
			logger.LogDbg(rsp.String())
		}
	}()
	return err
}

func (me *MyEtcd) lkey(key string) string {
	//return filepath.Join(me.prefix, key)
	return me.prefix + key + "/"
}

func (me *MyEtcd) Set(key, value string) bool {
	//lkey := filepath.Join(me.prefix, key)
	//lkey := me.prefix + key + "/"
	return set(me.Client, me.lkey(key), value)
}

func (me *MyEtcd) GrantSet(key, value string, ttl int) (clientv3.LeaseID, bool) {
	//lkey := filepath.Join(me.prefix, key)
	//lkey := me.prefix + key + "/"
	return grantSet(me.Client, me.lkey(key), value, ttl)
}

func (me *MyEtcd) Get(key string) ([]*EtcdKeyValue, bool) {
	//lkey := me.prefix + key + "/"
	return get(me.Client, me.lkey(key))
}

func (me *MyEtcd) GetByPrefix(key string) ([]*EtcdKeyValue, bool) {
	//lkey := me.prefix + key + "/"
	return getPrefix(me.Client, me.lkey(key))
}

func (me *MyEtcd) Watch(key string) WatchChan {
	rspchan := make(chan *WatchEvent)
	//lkey := me.prefix + key + "/"
	go watch(me.Client, me.lkey(key), rspchan)
	return rspchan
}
