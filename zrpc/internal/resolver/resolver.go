package resolver

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

const (
	// DirectScheme stands for direct scheme.
	DirectScheme = "direct"
	// DiscovScheme stands for discov scheme.
	DiscovScheme = "discov"
	// EndpointSepChar is the separator cha in endpoints.
	EndpointSepChar = ','

	subsetSize = 32
)

var (
	// EndpointSep is the separator string in endpoints.
	EndpointSep = fmt.Sprintf("%c", EndpointSepChar)

	dirBuilder directBuilder
	disBuilder discovBuilder
)

// RegisterResolver registers the direct and discov schemes to the resolver.
func RegisterResolver() {
	resolver.Register(&dirBuilder)
	// TODO 注册gRpc resolver Builder实现
	resolver.Register(&disBuilder)
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (r *nopResolver) Close() {
}

func (r *nopResolver) ResolveNow(options resolver.ResolveNowOptions) {
}
