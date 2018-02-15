package myetcd

// type EMonitorNode struct {
// 	PName string
// 	PAddr string
// }

// //1-monitor主要任务就是服务发现和注册;
// //	而这一切都依赖于Leader——通过选举产生Leader.
// type EMonitor struct {
// 	mine   EMonitorNode
// 	others map[string]EMonitorNode
// }

// func NewEMonitor() *EMonitor {
// 	return &EMonitor{
// 		EMonitorNode{
// 			progName(),
// 			progAddr(),
// 		},
// 		nil,
// 	}
// }
