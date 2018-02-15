package myetcd

import (
	"fmt"
	clientv3 "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	. "logger"
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
	_, err := me.Client.KeepAlive(context.TODO(), leaseid)
	return err
}

func (me *MyEtcd) lkey(key string) string {
	//return filepath.Join(me.prefix, key)
	return me.prefix + key + "/"
}

func (me *MyEtcd) Set(key, value string) bool {
	lkey := me.lkey(key)
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	LogDbg("put(%s) value(%s)", lkey, value)
	_, err := me.Client.Put(ctx, lkey, value)
	if err != nil {
		cancel()
		LogErr("Set %s", err.Error())
		return false
	}
	return true
}

func (me *MyEtcd) GrantSet(key, value string, ttl int) (clientv3.LeaseID, bool) {
	lkey := me.lkey(key)
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	resp, err := me.Client.Grant(context.TODO(), int64(ttl))
	if err != nil {
		LogErr("Grant error:%s", err.Error())
		return 0, false
	}
	LogDbg("grant put(%s) value(%s)", lkey, value)
	_, err = me.Client.Put(ctx, lkey, value, clientv3.WithLease(resp.ID))
	if err != nil {
		cancel()
		LogErr("Set %s", err.Error())
		return 0, false
	}
	return resp.ID, true
}

func (me *MyEtcd) Get(key string) ([]*EtcdKeyValue, bool) {
	lkey := me.lkey(key)
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	rsp, err := me.Client.Get(ctx, lkey)
	if err != nil {
		cancel()
		LogErr("Get %s", err.Error())
		return nil, false
	}
	result := make([]*EtcdKeyValue, len(rsp.Kvs))
	for idx, kv := range rsp.Kvs {
		result[idx] = &EtcdKeyValue{
			me.returnKey(kv.Key),
			kv.Value,
		}
	}
	return result, true
}

func (me *MyEtcd) GetByPrefix(key string) ([]*EtcdKeyValue, bool) {
	lkey := me.lkey(key)
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	LogDbg("getPrefix=%s", key)
	rsp, err := me.Client.Get(ctx, lkey, clientv3.WithPrefix())
	if err != nil {
		cancel()
		LogErr("Get %s", err.Error())
		return nil, false
	}
	result := make([]*EtcdKeyValue, len(rsp.Kvs))
	for idx, kv := range rsp.Kvs {
		result[idx] = &EtcdKeyValue{
			me.returnKey(kv.Key),
			kv.Value,
		}
	}
	return result, true
}

func (me *MyEtcd) returnKey(k EtcdKey) EtcdKey {
	return []byte(k)[len(me.prefix):]
}

func (me *MyEtcd) watch(key string, mychan chan<- *WatchEvent) {
	LogDbg("watch prefix key=%s", key)
	wchan := me.Client.Watch(context.TODO(), key, clientv3.WithPrefix())
	for {
		wrsp := <-wchan
		var we *WatchEvent
		for _, ev := range wrsp.Events {
			//LogDbg("watch : %+v", ev.Kv.String())
			switch int(ev.Type) {
			case 0: //mvccpb.PUT
				we = &WatchEvent{
					EE_PUT,
					EtcdKeyValue{
						me.returnKey(ev.Kv.Key),
						ev.Kv.Value,
					},
				}
			case 1: //mvccpb.DELETE
				we = &WatchEvent{
					EE_DEL,
					EtcdKeyValue{
						me.returnKey(ev.Kv.Key),
						ev.Kv.Value,
					},
				}
			default:
				LogErr("unknow type:%+v", ev.Type)
			}
			mychan <- we
		}
	}
	LogInfo("watch exit")
}

func (me *MyEtcd) Watch(key string) WatchChan {
	rspchan := make(chan *WatchEvent)
	//lkey := me.prefix + key + "/"
	go me.watch(me.lkey(key), rspchan)
	return rspchan
}
