package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/net/context"
	. "logger"
	"myetcd"
	"strings"
)

var (
	etcd    = flag.String("etcd", "192.168.13.100:2379", "address of etcd")
	sysname = flag.String("system", "knull", "name of system")
	key     = flag.String("key", "test", "prefix of key")
	op      = flag.String("op", "get", "operate")
	//
	eobj *myetcd.MyEtcd
)

func etcdinit() bool {
	addrs := strings.Split(*etcd, ",")
	eobj = myetcd.NewMyEtcd(*sysname, addrs)
	if eobj != nil {
		return (nil == eobj.Initialize())
	}
	return false
}

const endstr = "\n------------------------ E N D ---------------------\n"

func onEtcdGet() {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("onEtcdGet(%s)\n", *key))
	result, ok := eobj.GetByPrefix(*key)
	if ok {
		for idx, node := range result {
			buf.WriteString(fmt.Sprintf("\t%d : %s\n", idx, node.String()))
		}
	}

	buf.WriteString(endstr)
	LogDbg(buf.String())
}

func onEtcdSet() {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("onEtcdSet(%s)\n", *key))
	eobj.Set(*key, "just for test")
	eobj.Grant(context.TODO(), 3)
	// eobj.KeepAlive(5)
	buf.WriteString(endstr)
	LogDbg(buf.String())
}

func onEtcdGSet() {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("onEtcd-->G<--Set(%s)\n", *key))
	eobj.GrantSet(*key, "just for test", 3)
	// eobj.KeepAlive(5)
	buf.WriteString(endstr)
	LogDbg(buf.String())
}

func onWatch() {
	wchan := eobj.Watch(*key)
	for {
		rsp := <-wchan
		if rsp.Do == myetcd.EE_PUT {
			LogDbg("set value(%s) to key(%s)", rsp.EKV.Value, rsp.EKV.Key)
		} else {
			LogDbg("reomve value(%s) from key(%s)", rsp.EKV.Value, rsp.EKV.Key)
		}
	}
}

func EtcdRunning() {
	switch *op {
	case "gset":
		onEtcdGSet()
	case "set":
		onEtcdSet()
	case "watch":
		onWatch()
	// case "del":
	// 	onDel()
	default: //do get;
		onEtcdGet()
	}
}

func main() {
	flag.Parse()
	if !etcdinit() {
		return
	}
	EtcdRunning()
}
