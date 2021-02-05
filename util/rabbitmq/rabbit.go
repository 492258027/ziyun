package rabbitmq

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	_ "net/http/pprof"
	"time"
)

////////////////以下为创建链接(包括创建交换机，队列)，关闭链接， 生产者发布消息，消费者接收消息////////////////
type MqConn struct {
	Conn         *amqp.Connection
	Channel      *amqp.Channel
	Qqueue       *amqp.Queue
	exchangeName string // 交换机名称
	exchangeType string // 交换机类型
	queueName    string // 队列名称
	routingKey   string // key名称
}

func NewMqConn(amqpURI, exchangeName, exchangeType, queueName, routingKey string) *MqConn {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln("NewMqConn error : ", exchangeName, queueName, err)
		}
	}()

	//建立连接
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Errorln("Failed to connect to RabbitMQ:", err, amqpURI)
		return nil
	}

	//创建一个Channel
	channel, err := conn.Channel()
	if err != nil {
		log.Errorln("Failed to open a channel", err)
		return nil
	}

	//创建一个exchange
	err = channel.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		log.Errorln("Failed to declare a exchange", err)
		return nil
	}

	//创建一个queue
	qQueue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive 当Consumer关闭连接时，这个queue不被deleted
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Errorln("Failed to declare a queue", err)
		return nil
	}

	//绑定到exchange
	err = channel.QueueBind(
		qQueue.Name,  // name of the queue
		routingKey,   // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		log.Errorln("Failed to bind a queue", err)
		return nil
	}

	log.Infoln("NewMqConn:", conn.LocalAddr(), exchangeName, queueName, routingKey)

	return &MqConn{conn, channel, &qQueue, exchangeName, exchangeType, queueName, routingKey}
}

func (conn *MqConn) CloseMqConn() {
	if conn == nil {
		log.Errorln("CloseMqConn pointer nil")
		return
	}

	conn.Channel.Close()
	conn.Conn.Close()
}

func (conn *MqConn) Publish_mq(routingKey string, body []byte) error {
	if conn == nil {
		log.Errorln("Publish_mq pointer nil")
		return errors.New("Publish_mq pointer nil")
	}

	err := conn.Channel.Publish(
		conn.exchangeName, // exchange
		//conn.routingKey,   // routing key
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			DeliveryMode:    amqp.Persistent, //消息持久化，优先级低于队列持久化，exchange、队列、消息三者都持久化才有效
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
		})

	if err != nil {
		log.Errorln("mq publishing error:", err)
		return err
	} else {
		log.Infoln("publishing to mq :", len(body), body)
		return nil
	}
}

func (conn *MqConn) Consumer_mq(autoack bool) <-chan amqp.Delivery {
	if conn == nil {
		log.Errorln("Consumer_mq pointer = nil")
		return nil
	}

	//在手动回ack的模式下，设置consumer每次取消息的数量。自动回ack忽略此配置
	//	err := conn.Channel.Qos(
	//		1,     // prefetch count
	//		0,     // prefetch size
	//		false, // global
	//	)
	//	if err != nil {
	//		imlog.Log(imlog.LogLevelEnum_ERROR, "mq consume channel.Qos error", err)
	//	}

	msgch, err := conn.Channel.Consume(
		conn.Qqueue.Name, // queue
		"",               // consumer
		autoack,          // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)

	if err != nil {
		log.Errorln("mq consume error", err)
		return nil
	} else {
		return msgch
	}
}

////////////////以下为链接池相关， 构建链接池， 从链接池中获取可用链接， 归还链接////////////////
type MqPool struct {
	MqConnChan chan *MqConn
	//以下参数重建conn使用
	AmqpURI      string //URI
	ExchangeName string // 交换机名称
	ExchangeType string // 交换机类型
	QueueName    string // 队列名称
	RoutingKey   string // key名称
}

func NewMqPool(capacity int, amqpURI, exchangeName, exchangeType, queueName, routingKey string) *MqPool {
	if capacity <= 0 {
		log.Fatal("NewMqPool numConn error", capacity)
	}

	MqConnChan := make(chan *MqConn, capacity)

	for i := 0; i < capacity; i++ {
		conn := NewMqConn(amqpURI, exchangeName, exchangeType, queueName, routingKey)
		if conn == nil {
			log.Errorln("NewMqPool: NewMqConn create failure !!")
			return nil
		}

		select {
		case MqConnChan <- conn:
		default:
		}
	}
	return &MqPool{MqConnChan, amqpURI, exchangeName, exchangeType, queueName, routingKey}
}

func (mqPool *MqPool) OpenClient() *MqConn {
	if mqPool == nil {
		log.Fatal("OpenClient: mqPool pointer = nil")
	}

	select {
	case conn := <-mqPool.MqConnChan:
		return conn
	case <-time.After(time.Second * 1):
		return nil
	}
}

func (mqPool *MqPool) CloseClient(conn *MqConn) {
	if mqPool == nil {
		log.Fatal("CloseClient: mqPool pointer = nil")
	}

	if conn == nil {
		log.Errorln("CloseClient: MqConn pointer = nil")
		return
	}

	select {
	case mqPool.MqConnChan <- conn:
		return
	//此处default不可少，否则会阻塞。
	default:
		conn.CloseMqConn()
	}
}

////////////////以下对外暴露，初始化链接池，利于链接池与mq交互////////////////
//连接池不负责连接的可用性测试，
//使用者从连接池中取连接后，如果发现连接不可用，需要重新创建连接
//使用者放回的连接超出了连接池的容量，连接池会自动关闭链接，使用者只需往回放这个动作即可

var MqClient *MqPool

func InitRabbitMq(size int, amqpUri, exchangeName, exchangeType, queueName, routingKey string) {
	MqClient = NewMqPool(size, amqpUri, exchangeName, exchangeType, queueName, routingKey)
	if MqClient == nil {
		log.Fatal("InitRabbitMq: mqPool point nil!!")
	}
}

func PutDataToMq(mqPool *MqPool, routingKey string, body []byte) error {
	if mqPool == nil {
		log.Fatal("PutDataToMq: mqPool point nil!!")
	}

	//重试3次，不行重新创建链接
	for i := 0; i < 3; i++ {
		//从连接池获取连接
		conn := mqPool.OpenClient()
		//取链接成功
		if conn != nil {
			err := conn.Publish_mq(routingKey, body)
			// 发送失败
			if err != nil {
				log.Errorln("PutDataToMq: publish failure，close link!!  ", err)
				conn.CloseMqConn()
			} else { //发送成功, 放回连接池，return返回
				mqPool.CloseClient(conn)
				log.Infoln("PutDataToMq: get link from pool and publish success!! ")
				return nil
			}
		}
	}

	//要不就是取链接超时，要不就是取链接成功，但是发送超时。这两种情况都重新创建链接
	conn := NewMqConn(mqPool.AmqpURI, mqPool.ExchangeName, mqPool.ExchangeType, mqPool.QueueName, mqPool.RoutingKey)
	//如果创建链接失败， 直接return
	if conn == nil {
		log.Errorln("PutDataToMq: create new link failure!!")
		return errors.New(" create new link failure !!")
	} else {
		err := conn.Publish_mq(routingKey, body)
		// 发送失败
		if err != nil {
			log.Errorln("publish failure， close link!!  ", err)
			conn.CloseMqConn()
			return err
		} else { //发送成功, 放回连接池，return返回
			mqPool.CloseClient(conn)
			log.Infoln("PutDataToMq: create new link and publish success!!  ", err)
			return nil
		}
	}
}

//次函数取数据后只是打印，只是示例函数，实际中需要根据具体业务修改
func GetDataFromMq_t(mqPool *MqPool) {
	if mqPool == nil {
		log.Fatal("GetDataFromMq: mqPool point nil!!")
	}

	//从连接池获取连接。 只取一次， 取出来后链接不可用，直接重新创建，不再和链接池交互。
	conn := mqPool.OpenClient()

	defer func() {
		if err := recover(); err != nil {
			log.Errorln("GetDataFromMq: defer func, getDataFromMq error : ", err)
			conn.CloseMqConn()
			go GetDataFromMq_t(mqPool)
		}
	}()

HERE:
	//取链接成功
	if conn != nil {
		msgCh := conn.Consumer_mq(true)
		for {
			select {
			case hint, ok := <-msgCh:
				if ok {
					//解消息，并处理
					log.Infoln("GetDataFromMq: consume success", hint.Body)
				} else {
					log.Errorln("GetDataFromMq: get msg from chan failure! disconnect link and reconnect!")
					conn.CloseMqConn()
					// chan 读取失败， 需要隔一秒尝试创建一次，直到创建链接成功
					for {
						conn = NewMqConn(mqPool.AmqpURI, mqPool.ExchangeName, mqPool.ExchangeType, mqPool.QueueName, mqPool.RoutingKey)
						if conn != nil {
							log.Infoln("GetDataFromMq: create new link ok, goto here begin work!! chan func")
							goto HERE
						} else {
							log.Infoln("GetDataFromMq: create new link failure, sleep!! chan func")
							time.Sleep(time.Second)
						}
					}
				}
			}
		}
	}

	//pool 取链接失败，需要隔一秒尝试创建一次，直到创建链接成功
	for {
		log.Errorln("GetDataFromMq:  get link from pool failure!! pool func")
		conn = NewMqConn(mqPool.AmqpURI, mqPool.ExchangeName, mqPool.ExchangeType, mqPool.QueueName, mqPool.RoutingKey)
		if conn != nil {
			log.Infoln("GetDataFromMq: create new link ok, goto here begin work!! pool func")
			goto HERE
		} else {
			log.Infoln("GetDataFromMq:, create new link failure, sleep!! pool func")
			time.Sleep(time.Second)
		}
	}
}
