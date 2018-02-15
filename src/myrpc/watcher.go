package myrpc

import (
	"github.com/golang/glog"
	"google.golang.org/grpc/naming"
	"monitor"
	"myetcd"
	"sync/atomic"
)

// watcher is the implementaion of grpc.naming.Watcher
type Watcher struct {
	re            *resolver // re: Etcd Resolver
	isInitialized int32
}

// Close do nothing
func (w *Watcher) Close() {
}

func EtcdKV2GrpcUpdate(nodes []myetcd.EtcdKeyValue) []*naming.Update {
	dlen := len(nodes)
	updates := make([]*naming.Update, dlen)
	for i := 0; i < dlen; i++ {
		updates[i] = &naming.Update{Op: naming.Add, Addr: string(nodes[i].Value)}
		glog.Infof(">>>>>>>>init:%s", nodes[i].String())
	}
	return updates
}

// Next to return the updates
func (w *Watcher) Next() ([]*naming.Update, error) {

	// check if is initialized
	if atomic.CompareAndSwapInt32(&w.isInitialized, 0, 1) {
		result := monitor.GetNodes(w.re.serviceName)
		glog.Infof(">>>>>>>>watch(%s) init =%d", w.re.serviceName, len(result))
		return EtcdKV2GrpcUpdate(result), nil
	}

	// generate etcd Watcher
	rch := monitor.Watch(w.re.serviceName)
	for mrsp := range rch {
		if mrsp.EV == myetcd.EE_PUT {
			glog.Infof("Put : %s --> %s", mrsp.EKV.Key, mrsp.EKV.Value)
			return []*naming.Update{{Op: naming.Add, Addr: string(mrsp.EKV.Value)}}, nil
		}
		if mrsp.EV == myetcd.EE_DEL {
			glog.Infof("Del : %s --> %s", mrsp.EKV.Key, mrsp.EKV.Value)
			return []*naming.Update{{Op: naming.Delete, Addr: string(mrsp.EKV.Value)}}, nil
		}
	}
	return nil, nil
}

// func extractAddrs(resp *etcd3.GetResponse) []string {
// 	addrs := []string{}

// 	if resp == nil || resp.Kvs == nil {
// 		return addrs
// 	}

// 	for i := range resp.Kvs {
// 		if v := resp.Kvs[i].Value; v != nil {
// 			addrs = append(addrs, string(v))
// 		}
// 	}

// 	return addrs
// }
