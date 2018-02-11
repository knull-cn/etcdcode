package myetcd

// import (
// 	v3 "github.com/coreos/etcd/clientv3"
// 	. "logger"
// 	"time"
// )

// const (
// 	Election_Lease_Sec = 5
// 	split              = "/"
// )

// type MyInfoNode struct {
// 	name string
// 	addr string
// }

// type MyElection struct {
// 	key             string
// 	value           string
// 	mine            MyInfoNode
// 	leader          MyInfoNode
// 	etcdaddr        []string
// 	emtx            *EMutex
// 	v3cli           *v3.Client
// 	echan           chan bool
// 	electionrunning bool
// }

// func NewMyElection(addrs []string, name, addr string) *MyElection {
// 	cli, _ := NewClientObj(addrs)
// 	return &MyElection{
// 		key:      "/monitor/leader/election",
// 		value:    name + split + addr,
// 		etcdaddr: addrs,
// 		mine: MyInfoNode{
// 			name,
// 			addr,
// 		},
// 		emtx:            NewEMutex(addrs, "/monitor/leader"),
// 		v3cli:           cli,
// 		echan:           make(chan bool),
// 		electionrunning: false,
// 	}
// }

// func (me *MyElection) isLeader() bool {
// 	values, _ := get(me.v3cli, me.key)
// 	if len(values) > 0 {
// 		return string(values[0].Value) == me.value
// 	}
// 	return false
// }

// func (me *MyElection) GetLeader() {
// 	values, _ := get(me.v3cli, me.key)
// 	if len(values) > 0 {

// 	}
// }

// //选举;
// func (me *MyElection) Election() {
// 	LogDbg("%s Election-ing", me.mine.name)
// 	cli, _ := NewClientObj(me.etcdaddr)
// 	for {
// 		<-me.echan
// 		me.electionrunning = true
// 		if me.isLeader() {
// 			continue
// 		}
// 		me.emtx.Lock()
// 		set(cli, me.key, me.value)
// 		LogInfo("%s Campaign OK", me.mine.name)
// 	}
// }

// func (me *MyElection) GiveUp() {
// 	if me.isLeader() {
// 		LogInfo("%s GiveUp OK", me.mine.name)
// 		me.emtx.Unlock()
// 	}
// }

// func (me *MyElection) OnWatch(ee ETCD_EVENT, ekv *EtcdKeyValue) {
// 	if ee == EE_DEL {
// 		LogDbg("key(%s)value(%s) release", ekv.Key, ekv.Value)
// 		return
// 	}
// 	LogDbg("connect to key(%s)value(%s)")
// 	if string(ekv.Key) == me.key {
// 		if string(ekv.Value) != me.value {
// 			LogDbg("connect to key(%s)value(%s)")
// 		}
// 		me.echan <- true
// 	}
// }

// func (me *MyElection) Observe() {
// 	cli, _ := NewClientObj(me.etcdaddr)
// 	Watch(cli, me.key, me.OnWatch)
// }

// func (me *MyElection) Start() bool {
// 	// if !me.init() {
// 	// 	return false
// 	// }

// 	go me.Observe()
// 	go me.Election()
// 	time.Sleep(1 * time.Second)
// 	if !me.isLeader() {
// 		me.echan <- true
// 	}
// 	return true
// }
