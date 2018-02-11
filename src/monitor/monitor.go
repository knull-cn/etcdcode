package monitor

import (
	"bytes"
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
	mtx sync.RWMutex
	buf map[string]myetcd.EtcdValue
}

var mobj *Monitor

func (m *Monitor) GetNodes() (buf []myetcd.EtcdKeyValue) {
	// m.buf.Range(func(k, v interface{}) bool {
	// 	buf = append(buf, myetcd.EtcdKeyValue{
	// 		myetcd.EtcdKey(k.(string)),
	// 		v.(myetcd.EtcdValue),
	// 	})
	// 	//logger.LogDbg("GetNodes k=%+v;v=%+v;", k, v)
	// 	return true
	// })
	m.mtx.RLock()
	for k, v := range m.buf {
		buf = append(buf, myetcd.EtcdKeyValue{
			[]byte(k),
			v,
		})
	}
	m.mtx.RUnlock()
	//logger.LogInfo("push buf %d", len(buf))
	return buf
}

func (m *Monitor) onEvent(ev *myetcd.WatchEvent) {
	m.mtx.Lock()
	if ev.Do == myetcd.EE_DEL {
		delete(m.buf, string(ev.EKV.Key))
		//logger.LogInfo("delete %s", ev.EKV.Key)
	} else {
		m.buf[string(ev.EKV.Key)] = ev.EKV.Value
		//m.buf.Store(string(ev.EKV.Key), ev.EKV.Value)
		//logger.LogInfo("insert %s", ev.EKV.Key)
	}
	m.mtx.Unlock()
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

	//m.buf.Store(rkey, myetcd.EtcdValue(m.mine.PAddr))
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
		m.onEvent(rsp)
		rkey := m.registkey(m.mine.PName, m.mine.PAddr)
		if bytes.Contains(rsp.EKV.Key, []byte(rkey)) {
			if rsp.Do == myetcd.EE_DEL {
				logger.LogDbg("regist key(%s) deleted", rkey)
				m.etcd.Set(rkey, m.mine.PAddr)
				m.etcd.KeepAlive(MONITOR_TTL)
			}
		} else {
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

func (m *Monitor) initmap(cfg *MonitorCfg) bool {
	//get
	result, ok := m.etcd.GetByPrefix("monitor")
	if !ok {
		return false
	}
	m.mtx.Lock()
	for _, rnode := range result {
		for _, node := range cfg.Nodes {
			rkey := m.registkey(node.PName, node.PAddr)
			if bytes.Contains(rnode.Key, []byte(rkey)) {
				m.buf[string(rnode.Key)] = myetcd.EtcdValue(m.mine.PAddr)
				break
			}
		}
	}
	m.mtx.Unlock()
	return true
}

func newMonitor(cfg *MonitorCfg, etcdaddres []string) *Monitor {
	program := strings.Split(filepath.Base(os.Args[0]), ".")[0]
	logger.LogDbg("%s start monitor", program)
	var obj Monitor
	obj.buf = make(map[string]myetcd.EtcdValue)
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
	//map init;
	mobj.initmap(cfg)
	//
	go mobj.Register()

	return mobj.Monitor(), true
}

func GetNodes() []myetcd.EtcdKeyValue {
	return mobj.GetNodes()
}
