package protocol

import (
	"fmt"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"sync/atomic"
)

type MQServer struct {
	counter int64
}

func MQServerAttach(gsvr *grpc.Server) interface{} {
	server := &MQServer{}
	RegisterMQServiceServer(gsvr, server)
	return server
}

func (mqs *MQServer) HelloWorld(ctx context.Context, req *HelloWorldRequest) (*HelloWorldReply, error) {
	atomic.AddInt64(&mqs.counter, 1)
	glog.Info("HelloWorld ", req)
	return &HelloWorldReply{
		fmt.Sprintf("MQServer reply(%d)!", atomic.LoadInt64(&mqs.counter)),
	}, nil
}

//-----------------------------------------------
