package snowflake

// 雪花算法
// 把时间戳,工作机器ID, 序列号组合成一个64位int64
// 因为受限于sortedset中scores类型为float64,范围-2e53 ~ +2e53,实际可用53bit
// [12,51]这40位存放时间戳,[52,58]这7位存放机器id,[59,64]最后6位存放序列号
import (
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	machineID     int64 // 机器id占7位, 十进制范围是 [ 0, 127 ]
	sn            int64 // 序列号占5位,十进制范围是 [ 0, 63 ]
	lastTimeStamp int64 // 上次的时间戳(毫秒级) 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
)

const RefTime int64 = 1528236366000 // 2018年6月6日6时6分6秒时刻的毫秒值,所有时间戳减去这个数值,取相对时间差

var Snowflakelock sync.RWMutex

func InitSnowFlake(MillipedeID string) {
	lastTimeStamp = time.Now().UnixNano()/1e6 - RefTime

	var err error
	machineID, err = strconv.ParseInt(MillipedeID, 10, 64)

	if err != nil {
		log.Fatal("snowflake: machineID strconv error!")
	}

	// 把机器id左移6位,让出6位空间给序列号使用
	if machineID >= 0 && machineID <= 127 {
		machineID <<= 6
		log.Println("snowflake: machineID:", machineID)
	} else {
		machineID = 0
		log.Fatal("snowflake: macnineID invalid!")
	}
}

func GetSnowflakeId() int64 {
	//必须加锁,否则多个用户routine同时取消息id会出错
	Snowflakelock.Lock()
	defer Snowflakelock.Unlock()

	curTimeStamp := time.Now().UnixNano()/1e6 - RefTime

	// 同一毫秒，sn范围64，TPS为6.4万
	if curTimeStamp == lastTimeStamp {
		sn++
		// 序列号占5位,十进制范围是 [ 0, 63 ]
		if sn > 63 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano()/1e6 - RefTime
			lastTimeStamp = curTimeStamp
			sn = 0
		}

		// 机器id占用7位空间,序列号占用5位空间,所以左移12位;
		rightBinValue := curTimeStamp << 13

		id := rightBinValue | machineID | sn

		return id
	} else if curTimeStamp > lastTimeStamp {
		sn = 0
		lastTimeStamp = curTimeStamp
		rightBinValue := curTimeStamp << 13

		id := rightBinValue | machineID | sn

		return id
	} else {
		//理论上不会进入curTimeStamp < lastTimeStamp这个分支
		log.Println("curTimeStamp < lastTimeStamp, error!", curTimeStamp, lastTimeStamp)

		sn++
		curTimeStamp = lastTimeStamp
		rightBinValue := curTimeStamp << 13

		id := rightBinValue | machineID | sn

		return id
	}
}
