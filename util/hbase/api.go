package hbase

import (
	"errors"
	log "github.com/sirupsen/logrus"
	hb "ziyun/util/hbase/thrift/hbase"
)

func reconnect(hbaseClient *hb.THBaseServiceClient) error {
	hbaseClient.Transport.Close()
	log.Debugln("Hbase client reconnect")
	if err := hbaseClient.Transport.Open(); err != nil {
		log.Errorln("Error opening hbase socket", err.Error())
		return err
	}
	return nil
}

func Get(table []byte, get *hb.TGet) (r *hb.TResult_, err error) {
	hbaseClient := openClient()
	if hbaseClient == nil {
		log.Errorln("Get could not get hbase client from hbase connection pool")
		return nil, errors.New("Get could not get hbase client from hbase connection pool")
	}
	defer closeClient(hbaseClient)

	r, err = hbaseClient.Get(table, get)
	if err != nil {
		log.Warningln("Get hbase get error: ", err.Error())
		err = reconnect(hbaseClient)
		if err != nil {
			log.Errorln("Get reconnect hbase failed, error: ", err.Error())
			return nil, err
		}
		return hbaseClient.Get(table, get)
	}
	return
}

func GetMultiple(table []byte, gets []*hb.TGet) (r []*hb.TResult_, err error) {
	hbaseClient := openClient()
	if hbaseClient == nil {
		log.Errorln("GetMultiple could not get hbase client from hbase connection pool")
		return nil, errors.New("GetMultiple could not get hbase client from hbase connection pool")
	}
	defer closeClient(hbaseClient)

	r, err = hbaseClient.GetMultiple(table, gets)
	if err != nil {
		log.Warningln("GetMultiple hbase get error: ", err.Error())
		err = reconnect(hbaseClient)
		if err != nil {
			log.Errorln("GetMultiple reconnect hbase failed, error: ", err.Error())
			return nil, err
		}
		return hbaseClient.GetMultiple(table, gets)
	}
	return
}

func Put(table []byte, put *hb.TPut) (err error) {
	hbaseClient := openClient()
	if hbaseClient == nil {
		log.Errorln("Put could not get hbase client from hbase connection pool")
		return errors.New("Put could not get hbase client from hbase connection pool")
	}
	defer closeClient(hbaseClient)

	err = hbaseClient.Put(table, put)
	if err != nil {
		log.Warningln("Put hbase put error: ", err.Error())
		err = reconnect(hbaseClient)
		if err != nil {
			log.Errorln("Put reconnect hbase failed, error: ", err.Error())
			return err
		}
		return hbaseClient.Put(table, put)
	}
	return
}

func PutMultiple(table []byte, puts []*hb.TPut) (err error) {
	hbaseClient := openClient()
	if hbaseClient == nil {
		log.Errorln("PutMultiple could not get hbase client from hbase connection pool")
		return errors.New("PutMultiple could not get hbase client from hbase connection pool")
	}
	defer closeClient(hbaseClient)

	err = hbaseClient.PutMultiple(table, puts)
	if err != nil {
		log.Warningln("PutMultiple hbase put error: ", err.Error())
		err = reconnect(hbaseClient)
		if err != nil {
			log.Errorln("PutMultiple reconnect hbase failed, error: ", err.Error())
			return err
		}
		return hbaseClient.PutMultiple(table, puts)
	}
	return
}

func Del(table []byte, del *hb.TDelete) (err error) {
	hbaseClient := openClient()
	if hbaseClient == nil {
		log.Errorln("Del could not get hbase client from hbase connection pool")
		return errors.New("Del could not get hbase client from hbase connection pool")
	}
	defer closeClient(hbaseClient)

	err = hbaseClient.DeleteSingle(table, del)
	if err != nil {
		log.Warningln("Del hbase del error: ", err.Error())
		err = reconnect(hbaseClient)
		if err != nil {
			log.Errorln("Del reconnect hbase failed, error: ", err.Error())
			return err
		}
		return hbaseClient.DeleteSingle(table, del)
	}
	return
}

func DelMultiple(table []byte, dels []*hb.TDelete) (err error) {
	hbaseClient := openClient()
	if hbaseClient == nil {
		log.Errorln("DelMultiple could not get hbase client from hbase connection pool")
		return errors.New("DelMultiple could not get hbase client from hbase connection pool")
	}
	defer closeClient(hbaseClient)

	_, err = hbaseClient.DeleteMultiple(table, dels)
	if err != nil {
		log.Warningln("DelMultiple hbase del error: ", err.Error())
		err = reconnect(hbaseClient)
		if err != nil {
			log.Errorln("DelMultiple reconnect hbase failed, error: ", err.Error())
			return err
		}
		_, err = hbaseClient.DeleteMultiple(table, dels)
		return
	}
	return
}

func Scan(table []byte, scan *hb.TScan, numRows int32) (r []*hb.TResult_, err error) {
	hbaseClient := openClient()
	if hbaseClient == nil {
		log.Errorln("Scan could not get hbase client from hbase connection pool")
		return nil, errors.New("Scan could not get hbase client from hbase connection pool")
	}
	defer closeClient(hbaseClient)

	scannerId, err := hbaseClient.OpenScanner(table, scan)
	if err != nil {
		log.Warningln("Scan hbase open scan error: ", err.Error())
		reconnect(hbaseClient)
		scannerId, err = hbaseClient.OpenScanner(table, scan)
		if err != nil {
			log.Errorln("Scan open scanner error: ", err.Error())
			return nil, err
		}
	}
	defer hbaseClient.CloseScanner(scannerId)

	r, err = hbaseClient.GetScannerRows(scannerId, numRows)
	if err != nil {
		log.Errorln("Scan hbase get scan rows error: ", err.Error())
		return nil, err
	}
	return
}
