package monitor

import (
// "encoding/json"
// "os"
)

type MonitorNode struct {
	PName string `json:"pname"`
	PAddr string `json:"address"`
}

type MonitorCfg struct {
	SysName string        `json:"sysname"`
	Nodes   []MonitorNode `json:"mnodes"`
}
