package main

import (
	"bytes"
	"logger"
	"monitor"
	"time"
)

func MonitorTest() {
	cfg, err := parseConfig("./test.json")
	if err != nil {
		logger.LogErr("ParseConfig error:%s", err.Error())
		return
	}
	ok := monitor.StartMonitor(cfg, etcdaddr)
	if !ok {
		return
	}
	rspchan := monitor.Watch("")
	go func() {
		for {
			rsp := <-rspchan
			logger.LogDbg("onMonitor %+v", rsp)
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 7)
			var buf bytes.Buffer
			buf.WriteByte('\n')
			result := monitor.GetNodes("")
			for _, node := range result {
				buf.WriteByte('\t')
				buf.Write(node.Key)
				buf.WriteByte(' ')
				buf.Write(node.Value)
				buf.WriteByte('\n')
			}
			logger.LogDbg("GetNodes %s", buf.String())
		}
	}()
	//
	ch := make(chan bool)
	<-ch
	logger.LogDbg("MonitorTest exit")
}
