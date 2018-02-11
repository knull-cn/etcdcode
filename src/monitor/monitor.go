package monitor

import (
	"logger"
	"myetcd"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const (
	MONITOR_TTL = 5
)

type EMonitorNode struct {
	PName string
	PAddr string
}

type Monitor struct {
	mine     EMonitorNode
	monitors []EMonitorNode
	etcd     *myetcd.MyEtcd
	//
	buf sync.Map
}

var mobj *Monitor

func (m *Monitor) GetNodes() (buf []myetcd.EtcdKeyValue) {
	m.buf.Range(func(k, v interface{}) bool {
		buf = append(buf, myetcd.EtcdKeyValue{
			k.(myetcd.EtcdKey),
			v.(myetcd.EtcdValue),
		})
		return true
	})
	return buf
}

func (m *Monitor) onEvent(ev *myetcd.WatchEvent) {
	if ev.Do == myetcd.EE_DEL {
		m.buf.Delete(ev.EKV.Key)
	} else {
		m.buf.Store(ev.EKV.Key, ev.EKV.Value)
	}
}

func (m *Monitor) registkey(key, value string) string {
	//return filepath.Join("monitor", key, value)
	return "monitor" + "/" + key + "/" + value
}

func (m *Monitor) watchkey() string {
	//return filepath.Join("monitor")
	return "monitor"
}

func (m *Monitor) Register() {
	if m.mine.PName == "" {
		logger.LogInfo("not Register")
		return
	}
	logger.LogInfo("(%s) Register", m.mine.PName)
	rkey := m.registkey(m.mine.PName, m.mine.PAddr)
	leaseid, ok := m.etcd.GrantSet(rkey, m.mine.PAddr, MONITOR_TTL)
	if !ok {
		logger.LogErr("GrantSet failed")
		return
	}
	m.etcd.KeepAlive(leaseid)
	//
	// m.etcd.Watch(rkey, func(ev myetcd.ETCD_EVENT, _ *myetcd.EtcdKeyValue) {
	// 	if ev == myetcd.EE_DEL {
	// 		m.etcd.Set(rkey, m.mine.PAddr)
	// 		m.etcd.KeepAlive(MONITOR_TTL)
	// 	}
	// })
}

type MonitorRsp struct {
	EV  myetcd.ETCD_EVENT
	EKV *myetcd.EtcdKeyValue
}

type MonitorRspChan chan MonitorRsp

func (m *Monitor) Monitor() <-chan MonitorRsp {
	mrchan := make(chan MonitorRsp)
	go m.monitor(mrchan)
	return mrchan

}

func (m *Monitor) monitor(mrchan chan<- MonitorRsp) {
	rkey := m.watchkey()
	wchan := m.etcd.Watch(rkey)
	for {
		rsp := <-wchan
		rkey := m.registkey(m.mine.PName, m.mine.PAddr)
		if string(rsp.EKV.Key) == rkey {
			logger.LogDbg("regist key(%s) deleted", rkey)
			if rsp.Do == myetcd.EE_DEL {
				m.etcd.Set(rkey, m.mine.PAddr)
				m.etcd.KeepAlive(MONITOR_TTL)
			}
		} else {
			m.onEvent(rsp)
			mrchan <- MonitorRsp{
				rsp.Do,
				&myetcd.EtcdKeyValue{
					rsp.EKV.Key,
					rsp.EKV.Value,
				},
			}
		}
	}
}

func newMonitor(cfg *MonitorCfg, etcdaddres []string) *Monitor {
	program := strings.Split(filepath.Base(os.Args[0]), ".")[0]
	logger.LogDbg("%s start monitor", program)
	var obj Monitor
	obj.etcd = myetcd.NewMyEtcd(cfg.SysName, etcdaddres)
	for _, node := range cfg.Nodes {
		if node.PName == program {
			obj.mine.PName = node.PName
			obj.mine.PAddr = node.PAddr
		} else {
			obj.monitors = append(obj.monitors, EMonitorNode{
				node.PName,
				node.PAddr,
			})
		}
	}
	return &obj
}

func StartMonitor(path string, etcdaddres []string) (<-chan MonitorRsp, bool) {
	cfg, err := parseConfig(path)
	if err != nil {
		logger.LogErr("ParseConfig error:%s", err.Error())
		return nil, false
	}
	//
	mobj = newMonitor(cfg, etcdaddres)
	err = mobj.etcd.Initialize()
	if err != nil {
		logger.LogErr("Initialize error:%s", err.Error())
		return nil, false
	}
	//clean;

	//
	go mobj.Register()

	return mobj.Monitor(), true
}
