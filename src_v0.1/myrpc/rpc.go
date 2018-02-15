package myrpc

type RPCType int32

const (
	RPCT_NONE   RPCType = 0 //系统原生RPC;
	RPCT_GRPC           = 1 //google rpc;
	RPCT_THRIFT         = 2 //thrift
)
