package resolver

import (
	"strings"

	"github.com/tal-tech/go-zero/core/discov"
	"google.golang.org/grpc/resolver"
)

type discovBuilder struct{}

// TODO 这里实现了grpc的Builder接口，会被gRpc Client自动调用
func (d *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	hosts := strings.FieldsFunc(target.Authority, func(r rune) bool {
		return r == EndpointSepChar
	})
	// TODO 构建可以监听etcd上注册的服务变动信息的客户端
	sub, err := discov.NewSubscriber(hosts, target.Endpoint)
	if err != nil {
		return nil, err
	}

	update := func() {
		var addrs []resolver.Address
		// TODO subset会进行地址随机打散
		for _, val := range subset(sub.Values(), subsetSize) {
			addrs = append(addrs, resolver.Address{
				Addr: val,
			})
		}
		cc.UpdateState(resolver.State{
			// TODO 更新gRpc地址信息（所有接口的地址列表)
			Addresses: addrs,
		})
	}
	// TODO 添加etcd事件处理器
	sub.AddListener(update)
	update()

	return &nopResolver{cc: cc}, nil
}

func (d *discovBuilder) Scheme() string {
	return DiscovScheme
}
