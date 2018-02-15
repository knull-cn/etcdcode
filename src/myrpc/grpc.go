package myrpc

import (
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/keepalive"
	"net"
)

type ServerAttach func(s *grpc.Server, n string) interface{}

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
	svr := attach(grpcServer, info.Name)
	go grpcServer.Serve(lis)

	return svr, nil
}

//-----------------------------------------------
type GrpcClientInfo struct {
	Cred       GSinfoCred
	Name       string
	ServerName string
}

type ClientAttach func(s *grpc.ClientConn, n string) interface{}

func CreateClient(info GrpcClientInfo, attach ClientAttach) (interface{}, error) {
	//
	var opts []grpc.DialOption
	//
	resolver := NewResolver(info.ServerName)
	round := grpc.RoundRobin(resolver)
	opts = []grpc.DialOption{grpc.WithBalancer(round)}
	//
	cred := &info.Cred
	if cred.Pem != string("") && cred.Key != string("") {
		creds, err := credentials.NewServerTLSFromFile(cred.Pem, cred.Key)
		if err != nil {
			glog.Error("Failed to generate credentials ", err)
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(info.Name, opts...)
	if err != nil {
		glog.Error("Failed to generate credentials ", err)
		return nil, err
	}

	return attach(conn, info.Name), nil
}
