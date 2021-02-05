package hbase

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
	hb "ziyun/util/hbase/thrift/hbase"
)

var (
	hbaseChan chan *hb.THBaseServiceClient
)

func init() {
	hbaseChan = make(chan *hb.THBaseServiceClient, 10000)
}

func InitConnectionPool(cons int, addr string) bool {
	if cons <= 0 || addr == "" {
		return false
	}
	for i := 0; i < cons; i++ {
		protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			log.Errorln("error SplitHostPort hbase address:", err.Error())
			return false
		}
		transport, err := thrift.NewTSocketTimeout(net.JoinHostPort(host, port), time.Second*100)
		if err != nil {
			log.Errorln("error NewTSocketTimeout hbase address:", err.Error())
			return false
		}
		hbseClient := hb.NewTHBaseServiceClientFactory(transport, protocolFactory)
		if err := transport.Open(); err != nil {
			log.Errorln("Error opening socket", err.Error())
			return false
		}
		select {
		case hbaseChan <- hbseClient:
		default:
		}
	}

	return true
}

func openClient() *hb.THBaseServiceClient {
	select {
	case hbaseClient := <-hbaseChan:
		if !hbaseClient.Transport.IsOpen() {
			if err := hbaseClient.Transport.Open(); err != nil {
				log.Errorln("Error opening hbase socket", err.Error())
				return nil
			}
		}
		return hbaseClient
	case <-time.After(time.Second * 30):
		return nil
	}
	return nil
}

func closeClient(hbseClient *hb.THBaseServiceClient) {
	select {
	case hbaseChan <- hbseClient:
	default:
	}
}
