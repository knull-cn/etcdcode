package main

import (
	"logger"
	"monitor"
)

func MonitorTest() {
	rspchan, ok := monitor.StartMonitor("./test.json", etcdaddr)
	if !ok {
		return
	}
	go func() {
		for {
			rsp := <-rspchan
			logger.LogDbg("onMonitor %+v", rsp)
		}
	}()
	//
	ch := make(chan bool)
	<-ch
	logger.LogDbg("MonitorTest exit")
}
