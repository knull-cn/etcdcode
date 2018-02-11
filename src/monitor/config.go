package monitor

import (
	"encoding/json"
	"os"
)

type MonitorNode struct {
	PName string `json:"pname"`
	PAddr string `json:"address"`
}

type MonitorCfg struct {
	SysName string        `json:"sysname"`
	Nodes   []MonitorNode `json:"mnodes"`
}

func parseConfig(path string) (*MonitorCfg, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	//
	var cfg MonitorCfg
	decoder := json.NewDecoder(fp)
	err = decoder.Decode(&cfg)
	return &cfg, err
}
