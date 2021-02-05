package consul

import (
	"errors"
	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"log"
	"os"
	"strconv"
	"ziyun/util/bootstrap"
)

var (
	consulClient IConsulClient
	lb           ILoadBalance
	instanceId   string
)

var (
	NoInstanceExistedErr = errors.New("no instance exist in consul!")
	LoadBalanceErr       = errors.New("loadbalance choice error!")
)

func init() {
	var err error

	consulClient, err = NewConsulClient(bootstrap.ConsulConfig.Host, bootstrap.ConsulConfig.Port)
	if err != nil {
		log.Println("Get Consul Client failed")
		os.Exit(-1)
	}

	//注意接口需要指针赋值, 随机选择
	lb = &RandomLoadBalance{}

}

//根据服务名，随机选取agent
func Discover(serviceName string) (*api.AgentService, error) {
	instances := consulClient.DiscoverServices(serviceName)

	if len(instances) < 1 {
		log.Println("no available client for %s.", serviceName)
		return nil, NoInstanceExistedErr
	}

	// 组agent切片
	instanceList := make([]*api.AgentService, len(instances))
	for i := 0; i < len(instances); i++ {
		instanceList[i] = instances[i].(*api.AgentService)
	}

	// 随机的方式选取agent
	selectInstance, err := lb.SelectService(instanceList)
	if err != nil {
		return nil, LoadBalanceErr
	}

	return selectInstance, nil

}

//根据服务名和instanceId， 返回唯一的agent
func Discover_feature(serviceName string, feature string) (*api.AgentService, error) {
	instances := consulClient.DiscoverServices(serviceName)

	if len(instances) < 1 {
		log.Println("no available client for %s.", serviceName)
		return nil, NoInstanceExistedErr
	}

	// 找到指定特征的agent
	for i := 0; i < len(instances); i++ {
		p := instances[i].(*api.AgentService)
		if p.ID == feature {
			return p, nil
		}
	}

	return nil, NoInstanceExistedErr
}

func Register() {
	//// 实例失败，停止服务
	if consulClient == nil {
		panic(0)
	}

	//判空 instanceId,通过 go.uuid 获取一个服务实例ID
	instanceId = bootstrap.ConsulConfig.InstanceId

	if instanceId == "" {
		instanceId = bootstrap.ConsulConfig.ServiceName + uuid.NewV4().String()
	}

	//rpc的端口放入meta中
	meta := map[string]string{
		"rpcport": strconv.Itoa(bootstrap.RpcConfig.Port),
		"weight":  "1",
	}

	if !consulClient.Register(bootstrap.ConsulConfig.ServiceName, instanceId, "/health", bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port, meta) {
		log.Println("register service %s failed.", bootstrap.ConsulConfig.ServiceName)
		// 注册失败，服务启动失败
		panic(0)
	}
	log.Println(bootstrap.ConsulConfig.ServiceName+"-service for service %s success.", bootstrap.ConsulConfig.ServiceName)
}

func Deregister() {
	// 实例失败，停止服务
	if consulClient == nil {
		panic(0)
	}

	if !consulClient.DeRegister(instanceId) {
		log.Println("deregister for service %s failed.", bootstrap.ConsulConfig.ServiceName)
		panic(0)
	}
}
