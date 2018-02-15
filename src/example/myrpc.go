package main

import (
	"encoding/json"
	"logger"
	"monitor"
	"myrpc"
	"myrpc/protocol"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	gcfg    *monitor.MonitorCfg
	program string
)

func parseConfig(path string) (*monitor.MonitorCfg, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	//
	var cfg monitor.MonitorCfg
	decoder := json.NewDecoder(fp)
	err = decoder.Decode(&cfg)
	return &cfg, err
}

func StartMonitor() {
	cfg, err := parseConfig("./test.json")
	if err != nil {
		logger.LogErr("ParseConfig error:%s", err.Error())
		os.Exit(0)
		return
	}
	gcfg = cfg
	ok := monitor.StartMonitor(cfg, etcdaddr)
	if !ok {
		os.Exit(0)
		return
	}
	logger.LogDbg("StartMonitor running...")
}

func getnode(sname string) string {
	for _, node := range gcfg.Nodes {
		if node.PName == sname {
			return node.PAddr
		}
	}
	logger.LogErr("getnode(%s) not exist", sname)
	os.Exit(0)
	return ""
}

func StartMyRPCServer(sname string) {
	addr := getnode(sname)
	var info = myrpc.GrcpServerInfo{
		Address: addr,
		Name:    sname,
	}
	_, err := myrpc.CreateServer(info, protocol.MQServerAttach)
	if err != nil {
		logger.LogErr("CreateServer error:%s", err.Error())
		return
	}
	logger.LogDbg("StartMyRPCServer(%s) running", sname)
	ch := make(chan struct{})
	<-ch
}

func StartMyRPCClient(cname string, snames []string) {
	idx := int32(0)
	for n, node := range gcfg.Nodes {
		if node.PName == program {
			arr := strings.Split(node.PAddr, ":")
			if len(arr) == 2 {
				port, _ := strconv.ParseInt(arr[1], 10, 64)
				idx = int32(n*100000 + int(port))
			}
		}
	}
	for _, server := range snames {
		if server == program {
			continue
		}
		var info = myrpc.GrpcClientInfo{
			Name:       cname,
			ServerName: server,
		}
		obj, err := myrpc.CreateClient(info, protocol.MQClientAttach)
		if err != nil {
			logger.LogErr("CreateServer error:%s", err.Error())
			return
		}
		logger.LogDbg("CreateClient(%s) connect to (%s) ok", cname, server)
		go func() {
			cli := obj.(*protocol.MQClient)
			cli.SetIdx(int32(idx))
			for {
				rsp := cli.SendData("just for test")
				if rsp != "" {
					logger.LogDbg(">>>>>	%s", rsp)
				}
				time.Sleep(3 * time.Second)
			}
		}()
	}
}

func StartMyRPC() {
	program = filepath.Base(os.Args[0])
	logger.LogDbg("StartMyRPC(%s)", program)
	StartMonitor()
	go StartMyRPCServer(program)
	StartMyRPCClient(program, []string{"dbmgr", "room"})

	ch := make(chan struct{})
	<-ch
}
