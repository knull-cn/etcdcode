package myrpc

import (
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/keepalive"
	"net"
)

type ServerAttach func(s *grpc.Server) interface{}

type GSinfoCred struct {
	Pem string
	Key string
}

type GrcpServerInfo struct {
	Cred    GSinfoCred
	Address string
	Name    string
}

func CreateServer(info GrcpServerInfo, attach ServerAttach) (interface{}, error) {
	lis, err := net.Listen("tcp", info.Address)
	if err != nil {
		glog.Error("failed to listen ", err)
		return nil, err
	}
	var opts []grpc.ServerOption
	cred := &info.Cred
	if cred.Pem != string("") && cred.Key != string("") {
		creds, err := credentials.NewServerTLSFromFile(cred.Pem, cred.Key)
		if err != nil {
			glog.Error("Failed to generate credentials ", err)
			return nil, err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	svr := attach(grpcServer)
	go grpcServer.Serve(lis)

	return svr, nil
}
