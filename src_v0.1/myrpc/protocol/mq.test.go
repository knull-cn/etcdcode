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
	name    string
}

func MQServerAttach(gsvr *grpc.Server, name string) interface{} {
	server := &MQServer{}
	server.name = name
	RegisterMQServiceServer(gsvr, server)
	return server
}

func (mqs *MQServer) HelloWorld(ctx context.Context, req *HelloWorldRequest) (*HelloWorldReply, error) {
	atomic.AddInt64(&mqs.counter, 1)
	glog.Infof(" (%s)recive :%s", mqs.name, req.GetData())
	return &HelloWorldReply{
		fmt.Sprintf("MQServer reply(%d)!", atomic.LoadInt64(&mqs.counter)),
	}, nil
}

//-----------------------------------------------
var cid int32

func ClientIDX() int32 {
	return atomic.AddInt32(&cid, 1)
}

type MQClient struct {
	MQServiceClient
	counter int64
	idx     int32
	name    string
}

func (mc *MQClient) SetIdx(idx int32) {
	mc.idx = idx
}

func MQClientAttach(gcli *grpc.ClientConn, name string) interface{} {
	client := &MQClient{
		NewMQServiceClient(gcli),
		0,
		ClientIDX(),
		name,
	}
	return client
}

func (mc *MQClient) SendData(data string) string {
	atomic.AddInt64(&mc.counter, 1)
	req := HelloWorldRequest{
		fmt.Sprintf("client(%s:%d)send(%d) data=%s", mc.name, mc.idx, mc.counter, data),
	}
	reply, err := mc.HelloWorld(context.TODO(), &req)
	if err != nil {
		glog.Errorf("SendData error:%s", err.Error())
		return ""
	}
	return reply.Data
}
