package main

import (
	"flag"
	"fmt"
	"math/rand"
	"myetcd"
	"os"
	"strconv"
	"time"
)

var (
	//etcdaddr = []string{":2370", ":2371", ":2372", ":2379"}
	etcdaddr = []string{"192.168.116.251:2379"}
	baseport = int64(9000)
)

func giveuptest(enode *myetcd.MyElection) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		sec := rand.Intn(11)
		sec = 10
		time.Sleep(time.Second * time.Duration(sec))
		enode.GiveUp()
	}
}

func runner(idx int64) {
	vstr := strconv.FormatInt(int64(idx), 10)
	enode := myetcd.NewMyElection(etcdaddr, fmt.Sprintf("%8s", ""+vstr), ":"+strconv.FormatInt(int64(idx+baseport), 10))
	enode.Start()
	//go giveuptest(enode)
}

func etcdTest() {
	flag.Parse()
	if len(os.Args) > 1 {
		if os.Args[1] == "leader" {
			enode := myetcd.NewMyElection(etcdaddr, fmt.Sprintf("%8s", ""+"0"), ":"+strconv.FormatInt(int64(0+baseport), 10))
			//enode.GetLeader()
			enode.WatchLeader()
			return
		}
	}
	//
	runner(1)
	runner(202)
	runner(30303)
	runner(4040404)
	ch := make(chan bool)
	<-ch
}

func etcdTest1() {
	flag.Parse()
	port := int64(time.Now().Unix()%baseport + baseport)
	name := fmt.Sprintf("%d", port)
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	enode := myetcd.NewMyElection(etcdaddr, name, "127.0.0.1:"+strconv.FormatInt(port, 10))
	enode.Start()
	go giveuptest(enode)
	//
	ch := make(chan bool)
	<-ch
}
