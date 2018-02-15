package main

// import (
// 	"logger"
// 	"myetcd"
// )

// func WatchRunning() {
// 	myobj := myetcd.NewMyEtcd("monitor", etcdaddr)
// 	myobj.Initialize()
// 	myobj.Watch("monitor", func(ev myetcd.ETCD_EVENT, ekv *myetcd.EtcdKeyValue) {
// 		switch ev {
// 		case myetcd.EE_DEL:
// 			logger.LogDbg("delete %s %s", ekv.Key, ekv.Value)
// 		case myetcd.EE_PUT:
// 			logger.LogDbg("insert %s %s", ekv.Key, ekv.Value)
// 		}
// 	})
// }
