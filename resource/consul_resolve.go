package resource

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const (
	defaultPort = "8500"
)

var (
	errMissingAddr = errors.New("consul resolver: missing address")

	errAddrMisMatch = errors.New("consul resolver: invalied uri")

	errEndsWithColon = errors.New("consul resolver: missing port after port-separator colon")

	regexConsul, _ = regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z_]+)$")

	//单例模式
	builderInstance = &consulBuilder{}
)

func Init(host string, port int, token, sn string) {
	log.Println("calling consul init")
	//resolver.Register(CacheBuilder())
	resolver.Register(NewBuilder(host, port, token, sn))
}

type consulBuilder struct {
	host    string
	port    int
	token   string
	srvName string
}

type consulResolver struct {
	address              string
	token                string
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	name                 string
	disableServiceConfig bool
	Ch                   chan int
}

func NewBuilder(h string, p int, t, sn string) resolver.Builder {
	return &consulBuilder{host: h, port: p, token: t, srvName: sn}
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	cr := &consulResolver{
		address: fmt.Sprintf("%s:%d", cb.host, cb.port),
		token:   cb.token,
		name:    cb.srvName,
		cc:      cc,
		//disableServiceConfig: opts.DisableServiceConfig,
		Ch: make(chan int, 0),
	}
	go cr.watcher()
	return cr, nil

}

func (cr *consulResolver) watcher() {
	log.Printf("calling [%s] consul watcher", cr.name)
	config := api.DefaultConfig()
	config.Address = cr.address
	config.Token = cr.token
	client, err := api.NewClient(config)
	if err != nil {
		log.Printf("error create consul client: %v", err)
		return
	}
	t := time.NewTicker(2000 * time.Millisecond)
	defer func() {
		log.Println("watcher defer")
	}()
	i := 0
	for {
		select {
		case <-t.C:
			//fmt.Println("定时")
		case <-cr.Ch:
			//fmt.Println("ch call")
		}
		//api添加了 lastIndex   consul api中并不兼容附带lastIndex的查询
		services, _, err := client.Health().Service(cr.name, "", true, &api.QueryOptions{})
		if err != nil {
			log.Printf("error retrieving instances from Consul: %v", err)
			if i%5 == 0 {
				i = 0
			}
			i++
		} else {
			i = 0
		}
		newAddrs := make([]resolver.Address, 0)
		for _, service := range services {
			addr := net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))
			newAddrs = append(newAddrs, resolver.Address{
				Addr: addr,
				//type：不能是grpclb，grpclb在处理链接时会删除最后一个链接地址，不用设置即可 详见=> balancer_conn_wrappers => updateClientConnState
				ServerName: service.Service.Service,
			})
		}
		cr.cc.UpdateState(resolver.State{Addresses: newAddrs})
	}
}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	cr.Ch <- 1
}

func (cr *consulResolver) Close() {
}
