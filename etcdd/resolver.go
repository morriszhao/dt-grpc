package etcdd

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"log"
	"sync"
)

type etcdResolver struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cc         resolver.ClientConn
	etcdClient *clientv3.Client
	target     resolver.Target
	ipPool     sync.Map
}

func (er *etcdResolver) ResolveNow(resolver.ResolveNowOptions) {

	prefix := er.target.URL.Path
	spew.Dump(prefix)
	//获取etcd中 对应服务的 ip列表
	res, err := er.etcdClient.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Build etcd get addr failed; %v:", err)
		return
	}

	for _, kv := range res.Kvs {
		er.store(kv.Key, kv.Value)
	}

	er.updateState()
	go er.watcher()

}

func (er *etcdResolver) Close() {
	er.cancel()
}

//监听服务注册 是否有变化  以便能及时跟新ip_list
func (er *etcdResolver) watcher() {

	watchChan := er.etcdClient.Watch(context.Background(), "/"+er.target.URL.Path, clientv3.WithPrefix())

	for {
		select {
		case val := <-watchChan:
			for _, event := range val.Events {
				switch event.Type {
				case 0: // 0 是有数据增加
					er.store(event.Kv.Key, event.Kv.Value)
					log.Println("put:", string(event.Kv.Key))
					er.updateState()
				case 1: // 1是有数据减少
					log.Println("del:", string(event.Kv.Key))
					er.del(event.Kv.Key)
					er.updateState()
				}
			}
		case <-er.ctx.Done():
			return
		}

	}
}

func (er *etcdResolver) store(k, v []byte) {
	er.ipPool.Store(string(k), string(v))
}

func (er *etcdResolver) del(key []byte) {
	er.ipPool.Delete(string(key))
}

func (er *etcdResolver) updateState() {
	var addrList resolver.State
	er.ipPool.Range(func(k, v interface{}) bool {
		tA, ok := v.(string)
		if !ok {
			return false
		}
		addrList.Addresses = append(addrList.Addresses, resolver.Address{Addr: tA})
		return true
	})

	_ = er.cc.UpdateState(addrList)
}
