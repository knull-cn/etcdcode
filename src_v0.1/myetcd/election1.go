package myetcd

import (
	v3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"golang.org/x/net/context"
	. "logger"
	"strings"
	"sync"
	"time"
)

const (
	Election_Lease_Sec = 5
	split              = "/"
)

type MyInfoNode struct {
	name string
	addr string
}

type MyElection struct {
	obj       *concurrency.Election
	sess      *concurrency.Session
	key       string
	value     string
	mine      MyInfoNode
	leader    MyInfoNode
	etcdaddr  []string
	electchan chan bool
	//isLeader  bool
	mtx    sync.Mutex
	giveup chan int
}

func NewMyElection(addrs []string, name, addr string) *MyElection {
	return &MyElection{
		key:      "/monitor/leader",
		value:    name + split + addr,
		etcdaddr: addrs,
		mine: MyInfoNode{
			name,
			addr,
		},
		electchan: make(chan bool),
		giveup:    make(chan int),
	}
}

func (me *MyElection) init() bool {
	me.obj, me.sess = newElectionNode(me.key, me.etcdaddr)
	LogDbg("%s init ok", me.mine.name)
	return true
}

func newElectionNode(key string, etcdaddr []string) (*concurrency.Election, *concurrency.Session) {
	cli, ok := NewClientObj(etcdaddr)
	if !ok {
		return nil, nil
	}
	sess := newSession(cli)
	if sess == nil {
		cli.Close()
		return nil, nil
	}
	return concurrency.NewElection(sess, key), sess
}

func (me *MyElection) isLeader() bool {
	// if me.mine.name == me.leader.name {
	// 	return me.mine.addr == me.leader.addr
	// }
	//me.mtx.Lock()
	rsp, err := me.obj.Leader(context.TODO())
	//me.mtx.Unlock()
	if err != nil {
		return false
	}
	if len(rsp.Kvs) > 0 {
		//LogDbg("leader:%s", rsp.Kvs[0].String())
		return string(rsp.Kvs[0].Value) == me.value
	}
	return false
}

//选举;
func (me *MyElection) Election() {
	LogDbg("%s Election-ing", me.mine.name)
	for {
		<-me.electchan
		if me.isLeader() {
			continue
		}
		//me.mtx.Lock()
		err := me.obj.Campaign(context.TODO(), me.value)
		//me.mtx.Unlock()
		if err != nil {
			LogErr("%s Campaign error:%s", me.mine.name, err.Error())
			continue
		}
		if me.isLeader() {
			LogInfo("%s Campaign OK", me.mine.name)
			me.giveup <- 1
		} else {
			LogInfo("%s Campaign NOT OK", me.mine.name)
			me.obj.Proclaim(context.TODO(), me.value)
		}
	}
}

func (me *MyElection) GiveUp() {
	if !me.isLeader() {
		return
	}
	LogDbg("%s GiveUp", me.mine.name)
	//me.mtx.Lock()
	me.obj.Resign(context.TODO())
	//me.mtx.Unlock()
}

func (me *MyElection) Observe() {
	obj, _ := newElectionNode(me.key, me.etcdaddr)
	if obj == nil {
		return
	}
	//obj := me.obj
	time.Sleep(1 * time.Second)
	me.electchan <- true

	LogDbg("%s Observe-ing", me.mine.name)
	rspchan := obj.Observe(context.TODO())
	for {
		rsp := <-rspchan
		if len(rsp.Kvs) > 0 {
			kv := rsp.Kvs[0]
			//LogDbg("%s Observe : %s", me.mine.name, kv.String())
			//set me.leader;
			LogDbg("%s leader key=%s;value=%s", me.mine.name, kv.Key, kv.Value)
			lvalue := string(kv.Value)
			if lvalue != me.value {
				arr := strings.Split(lvalue, split)
				if len(arr) == 2 {
					me.leader.name = arr[0]
					me.leader.addr = arr[1]
				} else {
					LogErr("error formate lvalue=%s", lvalue)
				}
				me.electchan <- true
			}
		}
	}
}

func (me *MyElection) giveupTest() {
	for {
		<-me.giveup
		time.Sleep(time.Duration(3) * time.Second)
		me.GiveUp()
	}
}

func (me *MyElection) WatchLeader() {
	if !me.init() {
		LogDbg("key=%s;rev=%d;header=%s", me.obj.Key(), me.obj.Rev(), me.obj.Header().String())
		return
	}
	cli := me.sess.Client()
	wchan := cli.Watch(context.TODO(), me.key, v3.WithPrefix())
	LogDbg("WatchLeader key=%s", me.key)
	for {
		rsp := <-wchan
		for _, ev := range rsp.Events {
			if int(ev.Type) == 0 { //PUT;
				LogDbg("  PUT  leader :%s", ev.Kv.String())
			} else { //delete
				LogDbg("DELETE leader :%s", ev.Kv.String())
			}
		}
	}
}

func (me *MyElection) GetLeader() {
	if !me.init() {
		LogDbg("key=%s;rev=%d;header=%s", me.obj.Key(), me.obj.Rev(), me.obj.Header().String())
		return
	}
	rsp, err := me.obj.Leader(context.TODO())
	if err != nil {
		LogErr("GetLeader error=%s", err.Error())
		return
	}
	if len(rsp.Kvs) > 0 {
		kv := rsp.Kvs[0]
		LogDbg("GetLeader--->%s", kv.String())
	}
	LogDbg("key=%s;rev=%d;header=%s", me.obj.Key(), me.obj.Rev(), me.obj.Header().String())
	LogDbg("%s", rsp2show("election", me))
}

func (me *MyElection) Start() bool {
	if !me.init() {
		return false
	}

	go me.Observe()
	go me.Election()
	go me.giveupTest()
	return true
}
