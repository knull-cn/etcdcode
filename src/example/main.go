package main

import (
	"fmt"
	"sync"
)

func testsyncmap() {
	var tmap sync.Map
	tmap.Store("key", []byte("value"))
	tmap.Range(func(k, v interface{}) bool {
		fmt.Printf("k=%+v;v=%+v\n", k, v)
		return true
	})
}

func main() {
	//testsyncmap()
	//pageTest()
	//etcdTest()
	MonitorTest()
	//WatchRunning()
}
