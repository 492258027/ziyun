package log

import (
	"github.com/sirupsen/logrus"
)

//这个log给整个im包使用，方便统一修改日志格式和日志级别。
var Logrus *logrus.Logger

func init() {
	Logrus = logrus.New()
	//	一共7个级别日志输出： Trace << Debug << Info << Warning << Error << Fatal << Panic
	//	只输出不低于当前级别的日志数据, 日志的输出级别可以动态改变。
	Logrus.SetLevel(logrus.DebugLevel)

	//  内置了JSONFormatter和TextFormatter两种格式，也可以通过Formatter接口定义日志格式。 默认为TextFormatter
	//logrus.SetFormatter(&logrus.TextFormatter{
	//	ForceColors:               true,
	//	EnvironmentOverrideColors: true,
	//	TimestampFormat:           "2006-01-02 15:04:05", //时间格式
	//	// FullTimestamp:true,
	//	// DisableLevelTruncation:true,
	//})

	//重定向输出 默认为stderr
	//logrus.SetOutput(os.stderr)

	//日志定位： 定位函数及行号
	//Logrus.SetReportCaller(true)

	//hook机制，在logrus写入日志时拦截，可以添加字段，也可以把日志输出到es中
}
