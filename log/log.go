//Package log is a simple log system.
//Support query of timestamp
package log
import (
	"yar/transports"
)

type LogLevel int

const (

	LOG_DEBUG 	LogLevel = 0x01
	LOG_NOTICE 	LogLevel = 0x01 << 1
	LOG_NORMAL  LogLevel = 0x01 << 2
	LOG_WARNING	LogLevel = 0x01 << 3
	LOG_ERROR	LogLevel = 0x01 << 4
)

type Log interface {
	Append(conn transports.TransportConnection,level LogLevel,fmt string,params...interface{})
	//Query(level LogLevel,start time.Duration,end time.Duration)([]Log)
}

func LogParseLevel(level LogLevel) string{

	switch level {

	case LOG_DEBUG : {return "Debug"}
	case LOG_NOTICE: {return "Notice"}
	case LOG_NORMAL: {return "Normal"}
	case LOG_WARNING:{return "Warning"}
	case LOG_ERROR : {return "Error"}
	}

	return "Unknow"
}