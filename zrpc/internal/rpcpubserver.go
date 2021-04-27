package internal

import (
	"os"
	"strings"

	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/core/netx"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

// NewRpcPubServer returns a Server.
func NewRpcPubServer(etcdEndpoints []string, etcdKey, listenOn string, opts ...ServerOption) (Server, error) {
	// TODO discov/KeepAlive
	registerEtcd := func() error {
		// TODO 将要注册到Etcd的服务信息(host:port)
		pubListenOn := figureOutListenOn(listenOn)
		pubClient := discov.NewPublisher(etcdEndpoints, etcdKey, pubListenOn)
		return pubClient.KeepAlive()
	}
	server := keepAliveServer{
		registerEtcd: registerEtcd,
		Server:       NewRpcServer(listenOn, opts...),
	}

	return server, nil
}

type keepAliveServer struct {
	registerEtcd func() error
	Server
}

func (ags keepAliveServer) Start(fn RegisterFn) error {
	// TODO ags.registerEtcd()触发服务注册到etcd
	if err := ags.registerEtcd(); err != nil {
		return err
	}

	// TODO rpcServer.Start启动gRpc Server
	//  传入的RegisterFn参数用于将用户通过protobuf定义生成的gRpc服务描述信息注册到将要启动的gRpc Server
	//  好让gRpc Server能调用这些服务的实现处理调用方的请求
	return ags.Server.Start(fn)
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
