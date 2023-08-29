package etcdd

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

type etcdResolverBuilder struct {
	etcdClient *clientv3.Client
}

func NewEtcdResolverBuilder() *etcdResolverBuilder {

	//创建etcd 客户端链接
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatalf("client get etcd failed,error %v", err)
		return nil
	}

	return &etcdResolverBuilder{etcdClient: client}
}

func (erb *etcdResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	spew.Dump("build diaioyong")
	ctx, cancelFunc := context.WithCancel(context.Background())
	es := &etcdResolver{
		cc:         cc,
		etcdClient: erb.etcdClient,
		ctx:        ctx,
		cancel:     cancelFunc,
		target:     target,
	}

	es.ResolveNow(resolver.ResolveNowOptions{})
	return es, nil
}

func (erb *etcdResolverBuilder) Scheme() string {
	return "etcd"
}
