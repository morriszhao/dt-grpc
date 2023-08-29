package etcdd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

type etcdRegister struct {
	etcdCli       *clientv3.Client
	leaseId       clientv3.LeaseID
	ctx           context.Context
	cancel        context.CancelFunc
	keepAliveResp <-chan *clientv3.LeaseKeepAliveResponse
	conf          etcdConfig
}

func InitRegisterEtcd() {
	conf := initEtcdConfig()

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{conf.host},
		DialTimeout: time.Duration(conf.timeout) * time.Second,
	})

	if err != nil {
		log.Fatalf("new etcd client failed,error %v \n", err)
		return
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	e := &etcdRegister{
		etcdCli: client,
		ctx:     ctx,
		cancel:  cancelFunc,
		conf:    conf,
	}

	e.createLease()
	e.bindLease()
	e.keepAlive()
	go e.watcher()

}

//创建租约
func (e *etcdRegister) createLease() {
	res, err := e.etcdCli.Grant(e.ctx, int64(e.conf.expired))
	if err != nil {
		log.Fatalf("create Lease failed,error %v \n", err)
	}

	e.leaseId = res.ID
}

//绑定租约
func (e *etcdRegister) bindLease() {
	put, err := e.etcdCli.Put(e.ctx, e.conf.serverName, e.conf.serverPath, clientv3.WithLease(e.leaseId))
	if err != nil {
		log.Fatalf("bindLease failed,error %v \n", err)
	}
	log.Printf("bindlease success! %v \n", put)
}

//续租 发送心跳 表明服务正常
func (e *etcdRegister) keepAlive() {
	alive, err := e.etcdCli.KeepAlive(e.ctx, e.leaseId)
	if err != nil {
		log.Fatalf("keepAlive failed,error %v \n", err)
		return
	}

	e.keepAliveResp = alive
}

func (e *etcdRegister) watcher() {
	for {
		select {
		case l := <-e.keepAliveResp:
			log.Printf("续约成功,val:%+v \n", l)
		case <-e.ctx.Done():
			log.Printf("续约关闭")
			return
		}
	}
}

func (e *etcdRegister) close() {
	e.cancel()
	log.Printf("closed...\n")
	e.etcdCli.Revoke(e.ctx, e.leaseId)
	e.etcdCli.Close()
}

type etcdConfig struct {
	serverName string //注册到etcd 的服务名称
	serverPath string //注册到etcd 的服务地址
	host       string //etcd 本身地址
	timeout    int    //超时设置
	expired    int    //租约超时
	keepalive  int    //心跳间隔
}

func initEtcdConfig() etcdConfig {

	//todo  从配置文件读取
	return etcdConfig{
		serverName: "/grpc/user_service/1",
		serverPath: "localhost:8888",
		host:       "127.0.0.1:2379",
		timeout:    1,
		expired:    60,
		keepalive:  10,
	}
}
