package main

import (
	"bytes"
	"logger"
	"monitor"
	"time"
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

	go func() {
		for {
			time.Sleep(time.Second * 7)
			var buf bytes.Buffer
			buf.WriteByte('\n')
			result := monitor.GetNodes()
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
