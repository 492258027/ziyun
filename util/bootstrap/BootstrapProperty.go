package bootstrap

var (
	HttpConfig    HttpConf
	RpcConfig     RpcConf
	ConsulConfig  ConsulConf
	HystrixConfig HystrixConf
	GatewayConfig GatewayConf
	JwtConfig     JwtConf
	Wsconfig      WsConf
	Mgrconfig     MgrConf
	Mqconfig      MqConf
	Redisconfig   RedisConf
	Mysqlconfig   MysqlConf
)

//Http配置
type HttpConf struct {
	Host string
	Port int
}

// RPC配置
type RpcConf struct {
	Host string
	Port int
}

//服务发现与注册配置
type ConsulConf struct {
	Host        string
	Port        int
	ServiceName string
	InstanceId  string
	Weight      int
}

// hystrix界面
type HystrixConf struct {
	Host string
	Port int
}

type GatewayConf struct {
	ConsulAuthName     string
	ConsulOpStringName string
}

// jwt配置
type JwtConf struct {
	JwtSecretKey string
}

// websocket
type WsConf struct {
	Host string
	Port int
}

// im的管理中心
type MgrConf struct {
	Host string
	Port int
}

// rabbitmq
type MqConf struct {
	PoolSize     int
	AmqpURI      string
	ExchangeName string
	ExchangeType string
	QueueName    string
	RoutingKey   string
}

//redis
type RedisConf struct {
	ClusterIPs   []string
	PoolSize     int
	MinIdleConns int
	Password     string
}

//mysql
type MysqlConf struct {
	DbUser          string
	DbPassword      string
	DbHost          string
	DbPort          int
	DbName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
}
