package myetcd

import (
	"encoding/json"
	"fmt"
)

func rsp2show(name string, data interface{}) string {
	jdata, _ := json.Marshal(data)
	return fmt.Sprintf("%s : %s", name, jdata)
}
