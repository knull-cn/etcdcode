package myetcd

import (
	v3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"golang.org/x/net/context"
	. "logger"
)

// func NewClusterCli(cli *v3.Client) {
// 	cluster := v3.NewCluster(cli)
// 	//v3.NewClusterFromClusterClient()
// }

func NewClientObj(EtcdAddr []string) (*v3.Client, bool) {
	cli, err := v3.New(v3.Config{
		Endpoints:   EtcdAddr,
		DialTimeout: timeout,
	})
	if err != nil {
		LogErr("NewClientObj failed:%s", err.Error())
		return nil, false
	}
	return cli, true
}

func grantSet(cli *v3.Client, key, value string, ttl int) (v3.LeaseID, bool) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	resp, err := cli.Grant(context.TODO(), int64(ttl))
	if err != nil {
		LogErr("Grant error:%s", err.Error())
		return 0, false
	}
	LogDbg("put(%s) value(%s)", key, value)
	_, err = cli.Put(ctx, key, value, v3.WithLease(resp.ID))
	if err != nil {
		cancel()
		LogErr("Set %s", err.Error())
		return 0, false
	}
	return resp.ID, true
}

func set(cli *v3.Client, key, value string) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	LogDbg("put(%s) value(%s)", key, value)
	_, err := cli.Put(ctx, key, value)
	if err != nil {
		cancel()
		LogErr("Set %s", err.Error())
		return false
	}
	return true
}

func del(cli *v3.Client, key string) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	_, err := cli.Delete(ctx, key)
	if err != nil {
		cancel()
		LogErr("Del %s", err.Error())
		return false
	}
	return true
}

func get(cli *v3.Client, key string) ([]*EtcdKeyValue, bool) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	rsp, err := cli.Get(ctx, key)
	if err != nil {
		cancel()
		LogErr("Get %s", err.Error())
		return nil, false
	}
	result := make([]*EtcdKeyValue, len(rsp.Kvs))
	for idx, kv := range rsp.Kvs {
		result[idx] = &EtcdKeyValue{
			kv.Key,
			kv.Value,
		}
	}
	return result, true
}

func getPrefix(cli *v3.Client, key string) ([]*EtcdKeyValue, bool) {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	LogDbg("getPrefix=%s", key)
	rsp, err := cli.Get(ctx, key, v3.WithPrefix())
	if err != nil {
		cancel()
		LogErr("Get %s", err.Error())
		return nil, false
	}
	result := make([]*EtcdKeyValue, len(rsp.Kvs))
	for idx, kv := range rsp.Kvs {
		result[idx] = &EtcdKeyValue{
			kv.Key,
			kv.Value,
		}
	}
	return result, true
}

// type WachCallback func(ETCD_EVENT, *EtcdKeyValue)

// func Watch(cli *v3.Client, key string, cb WachCallback) {
// 	LogDbg("watch prefix key=%s", key)
// 	wchan := cli.Watch(context.TODO(), key, v3.WithPrefix())
// 	for {
// 		wrsp := <-wchan
// 		// if !wrsp.Canceled {
// 		// 	continue
// 		// }
// 		for _, ev := range wrsp.Events {
// 			LogDbg(ev.Kv.String())
// 			switch int(ev.Type) {
// 			case 0: //mvccpb.PUT
// 				cb(EE_PUT, &EtcdKeyValue{
// 					ev.Kv.Key,
// 					ev.Kv.Value,
// 				})
// 			case 1: //mvccpb.DELETE
// 				cb(EE_DEL, &EtcdKeyValue{
// 					ev.Kv.Key,
// 					ev.Kv.Value,
// 				})
// 			default:
// 				LogErr("unknow type:%+v", ev.Type)
// 			}
// 		}
// 	}
// 	LogInfo("watch exit")
// }

type WatchEvent struct {
	Do  ETCD_EVENT
	EKV EtcdKeyValue
}

type WatchChan <-chan *WatchEvent

func watch(cli *v3.Client, key string, mychan chan<- *WatchEvent) {
	LogDbg("watch prefix key=%s", key)
	wchan := cli.Watch(context.TODO(), key, v3.WithPrefix())
	for {
		wrsp := <-wchan
		var we *WatchEvent
		for _, ev := range wrsp.Events {
			LogDbg("watch : %+v", ev.Kv.String())
			switch int(ev.Type) {
			case 0: //mvccpb.PUT
				we = &WatchEvent{
					EE_PUT,
					EtcdKeyValue{
						ev.Kv.Key,
						ev.Kv.Value,
					},
				}
			case 1: //mvccpb.DELETE
				we = &WatchEvent{
					EE_DEL,
					EtcdKeyValue{
						ev.Kv.Key,
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

func newSession(cli *v3.Client) *concurrency.Session {
	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	// lease := v3.NewLease(cli)
	// lease.Grant()
	rsp, err := cli.Grant(ctx, Election_Lease_Sec)
	if err != nil {
		cancel()
		LogErr("Grant %s", err.Error())
		return nil
	}
	session, err1 := concurrency.NewSession(cli, concurrency.WithLease(rsp.ID))
	if err1 != nil {
		cancel()
		LogErr("NewSession %s", err1.Error())
		return nil
	}
	return session
}

// func Election(key string) *concurrency.Election {
// 	session := newSession()
// 	if session != nil {
// 		return nil
// 	}
// 	//
// 	return concurrency.NewElection(session, key)
// }
