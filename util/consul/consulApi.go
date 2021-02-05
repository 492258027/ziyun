package consul

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"strconv"
	"sync"
)

type IConsulClient interface {

	/**
	 * 服务注册接口
	 * @param serviceName 服务名
	 * @param instanceId 服务实例Id
	 * @param instancePort 服务实例端口
	 * @param healthCheckUrl 健康检查地址
	 * @param instanceHost 服务实例地址
	 * @param meta 服务实例元数据
	 */
	Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string) bool

	/**
	 * 服务注销接口
	 * @param instanceId 服务实例Id
	 */
	DeRegister(instanceId string) bool

	/**
	 * 发现服务实例接口
	 * @param serviceName 服务名
	 */
	DiscoverServices(serviceName string) []interface{}
}

type ConsulClient struct {
	//consul 的host和port
	Host string
	Port int

	//consul 的config 和 client
	config *api.Config
	client consul.Client

	//锁及服务实例缓存字段
	mutex        sync.Mutex
	instancesMap sync.Map
}

func NewConsulClient(consulHost string, consulPort int) (*ConsulClient, error) {
	//consul 原生库的config和client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(consulPort)
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}

	//go kit 的consul库
	client := consul.NewClient(apiClient)
	return &ConsulClient{
		Host:   consulHost,
		Port:   consulPort,
		config: consulConfig,
		client: client,
	}, err
}

func (c *ConsulClient) Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string) bool {

	// 1. 构建服务实例元数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta:    meta,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}

	// 2. 发送服务注册到 Consul 中
	err := c.client.Register(serviceRegistration)

	if err != nil {
		log.Println("Register Service Error!")
		return false
	}
	log.Println("Register Service Success!")
	return true
}

func (c *ConsulClient) DeRegister(instanceId string) bool {

	// 构建包含服务实例 ID 的元数据结构体
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId,
	}
	// 发送服务注销请求
	err := c.client.Deregister(serviceRegistration)

	if err != nil {
		log.Println("Deregister Service Error!")
		return false
	}
	log.Println("Deregister Service Success!")

	return true
}

func (c *ConsulClient) DiscoverServices(serviceName string) []interface{} {
	//  该服务已监控并缓存
	instanceList, ok := c.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	}
	// 申请锁
	c.mutex.Lock()
	// 再次检查是否监控
	instanceList, ok = c.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]interface{})
	} else {
		// 注册监控
		go func() {
			// 使用 consul 服务实例监控来监控某个服务名的服务实例列表变化
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return // 数据异常，忽略
				}
				// 没有服务实例在线
				if len(v) == 0 {
					c.instancesMap.Store(serviceName, []interface{}{})
				}
				var healthServices []interface{}
				for _, service := range v {
					if service.Checks.AggregatedStatus() == api.HealthPassing {
						healthServices = append(healthServices, service.Service)
					}
				}
				c.instancesMap.Store(serviceName, healthServices)
			}
			defer plan.Stop()
			plan.Run(c.config.Address)
		}()
	}
	defer c.mutex.Unlock()

	// 根据服务名请求服务实例列表
	entries, _, err := c.client.Service(serviceName, "", false, nil)
	if err != nil {
		c.instancesMap.Store(serviceName, []interface{}{})
		log.Println("Discover Service Error!")
		return nil
	}
	instances := make([]interface{}, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = entries[i].Service
	}
	c.instancesMap.Store(serviceName, instances)
	return instances
}
