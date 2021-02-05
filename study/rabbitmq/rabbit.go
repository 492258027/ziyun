package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	_ "net/http/pprof"
	"time"
)

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
			log.Error("NewMqConn error : ", exchangeName, queueName, err)
		}
	}()

	//建立连接
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Error("Failed to connect to RabbitMQ:", err, amqpURI)
		return nil
	}

	//创建一个Channel
	channel, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a channel", err)
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
		log.Error("Failed to declare a exchange", err)
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
		log.Error("Failed to declare a queue", err)
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
		log.Error("Failed to bind a queue", err)
		return nil
	}

	log.Info("NewMqConn:", " addr: ", conn.LocalAddr(), " exchangeName: ", exchangeName, " queueName: ", queueName, " routingKey: ", routingKey)

	return &MqConn{conn, channel, &qQueue, exchangeName, exchangeType, queueName, routingKey}
}

func (conn *MqConn) CloseMqConn() {
	if conn == nil {
		log.Error("CloseMqConn pointer nil")
		return
	}

	conn.Channel.Close()
	conn.Conn.Close()
}

func (conn *MqConn) Publish_mq(body []byte) error {
	if conn == nil {
		log.Error("Publish_mq pointer nil")
		return errors.New("Publish_mq pointer nil")
	}

	err := conn.Channel.Publish(
		conn.exchangeName, // exchange
		conn.routingKey,   // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			DeliveryMode:    amqp.Persistent, //消息持久化，优先级低于队列持久化，exchange、队列、消息三者都持久化才有效
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
		})

	if err != nil {
		log.Error("mq publishing error:", err)
		return err
	} else {
		log.Info("publishing to mq :", len(body), body)
		return nil
	}
}

func (conn *MqConn) Consumer_mq(autoack bool) <-chan amqp.Delivery {
	if conn == nil {
		log.Error("Consumer_mq pointer = nil")
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
		log.Error("mq consume error", err)
		return nil
	} else {
		return msgch
	}
}

type MqPool struct {
	MqConnChan chan *MqConn
	//以下参数重建conn使用
	amqpURI      string //URI
	exchangeName string // 交换机名称
	exchangeType string // 交换机类型
	queueName    string // 队列名称
	routingKey   string // key名称
}

func NewMqPool(capacity int, amqpURI, exchangeName, exchangeType, queueName, routingKey string) *MqPool {
	if capacity <= 0 {
		log.Fatal("NewMqPool numConn error", capacity)
	}

	MqConnChan := make(chan *MqConn, capacity)

	for i := 0; i < capacity; i++ {
		conn := NewMqConn(amqpURI, exchangeName, exchangeType, queueName, routingKey)
		if conn == nil {
			log.Error("NewMqPool: NewMqConn create failure !!")
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
		log.Error("CloseClient: MqConn pointer = nil")
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

//连接池不负责连接的可用性测试，
//使用者从连接池中取连接后，如果发现连接不可用，需要重新创建连接
//使用者放回的连接超出了连接池的容量，连接池会自动关闭链接，使用者只需往回放这个动作即可
// delay ???? 如果mq队列被删除, Publish_mq也返回成功

func PutDataToMq(mqPool *MqPool, body []byte) error {
	if mqPool == nil {
		log.Fatal("PutDataToMq: mqPool point nil!!")
	}

	//重试3次，不行重新创建链接
	for i := 0; i < 3; i++ {
		//从连接池获取连接
		conn := mqPool.OpenClient()
		//取链接成功
		if conn != nil {
			err := conn.Publish_mq(body)
			// 发送失败
			if err != nil {
				log.Error("PutDataToMq: publish failure，close link!!  ", err)
				conn.CloseMqConn()
			} else { //发送成功, 放回连接池，return返回
				mqPool.CloseClient(conn)
				log.Info("PutDataToMq: get link from pool and publish success!! ")
				return nil
			}
		}
	}

	//要不就是取链接超时，要不就是取链接成功，但是发送超时。这两种情况都重新创建链接
	conn := NewMqConn(mqPool.amqpURI, mqPool.exchangeName, mqPool.exchangeType, mqPool.queueName, mqPool.routingKey)
	//如果创建链接失败， 直接return
	if conn == nil {
		log.Error("PutDataToMq: create new link failure!!")
		return errors.New(" create new link failure !!")
	} else {
		err := conn.Publish_mq(body)
		// 发送失败
		if err != nil {
			log.Error("publish failure， close link!!  ", err)
			conn.CloseMqConn()
			return err
		} else { //发送成功, 放回连接池，return返回
			mqPool.CloseClient(conn)
			log.Info("PutDataToMq: create new link and publish success!!  ", err)
			return nil
		}
	}
}

func GetDataFromMq(mqPool *MqPool) {
	if mqPool == nil {
		log.Fatal("GetDataFromMq: mqPool point nil!!")
	}

	//从连接池获取连接。 只取一次， 取出来后链接不可用，直接重新创建，不再和链接池交互。
	conn := mqPool.OpenClient()

	defer func() {
		if err := recover(); err != nil {
			log.Error("GetDataFromMq: defer func, getDataFromMq error : ", err)
			conn.CloseMqConn()
			go GetDataFromMq(mqPool)
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
					log.Info("GetDataFromMq: consume success", hint.Body)
				} else {
					log.Error("GetDataFromMq: get msg from chan failure! disconnect link and reconnect!")
					conn.CloseMqConn()
					// chan 读取失败， 需要隔一秒尝试创建一次，直到创建链接成功
					for {
						conn = NewMqConn(mqPool.amqpURI, mqPool.exchangeName, mqPool.exchangeType, mqPool.queueName, mqPool.routingKey)
						if conn != nil {
							log.Info("GetDataFromMq:,  create new link ok, goto here begin work!! chan func")
							goto HERE
						} else {
							log.Info("GetDataFromMq:, create new link failure, sleep!! chan func")
							time.Sleep(time.Second)
						}
					}
				}
			}
		}
	}

	//pool 取链接失败，需要隔一秒尝试创建一次，直到创建链接成功
	for {
		log.Error("GetDataFromMq:  get link from pool failure!! pool func")
		conn = NewMqConn(mqPool.amqpURI, mqPool.exchangeName, mqPool.exchangeType, mqPool.queueName, mqPool.routingKey)
		if conn != nil {
			log.Info("GetDataFromMq: create new link ok, goto here begin work!! pool func")
			goto HERE
		} else {
			log.Info("GetDataFromMq:, create new link failure, sleep!! pool func")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	uri := "amqp://root:root@192.168.73.3:5672/"
	pp := NewMqPool(2, uri, "exchange_xu", "direct", "queue_xu", "xu")

	go GetDataFromMq(pp)

	for {
		PutDataToMq(pp, []byte("hello boy!"))
		time.Sleep(1 * time.Second)
	}

	ch := make(chan int)
	<-ch
}

/*
func PutDataToMq2(mqPool *MqPool, body []byte) error {

		//从连接池获取连接
		conn := mqPool.OpenClient()
		//取链接成功
		if conn != nil {
			err := conn.Publish_mq(body)
			// 发送失败
			if err !=nil{
				log.Error("publish failure!!  ", err)
				conn.CloseMqConn()
				goto HERE
			}else{ //发送成功
				mqPool.CloseClient(conn)
				return nil
			}
		}else{	//取链接超时
			log.Error("get link from pool timeout!!")
			goto HERE
		}

		//重新创建链接，重试put操作.
		HERE:
			conn = NewMqConn(mqPool.amqpURI, mqPool.exchangeName, mqPool.exchangeType, mqPool.queueName, mqPool.routingKey)
			//如果创建链接失败， 直接return
			if conn == nil {
				log.Error("PutDataToMq: NewMqConn create failure!!")
				return errors.New("NewMqConn create failure !!")
			}else{
				err := conn.Publish_mq(body)
				// 发送失败
				if err !=nil{
					log.Error("publish failure!!  ", err)
					conn.CloseMqConn()
					return err
				}else{ //发送成功
					mqPool.CloseClient(conn)
					return nil
				}
			}
}
*/
