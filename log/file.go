package log
import (
	"os"
	"yar/transports"
	"fmt"
)

type FileLog struct {

	file *os.File

}

func NewFileLog(path string)(l *FileLog,err error){

	l = new(FileLog)
	f,err := os.OpenFile(path,os.O_CREATE | os.O_RDWR | os.O_APPEND,os.ModePerm)

	if err != nil {
		return nil,err
	}

	l.file = f
	return l,nil
}

func (log *FileLog)Append(conn transports.TransportConnection,level LogLevel,f string,params...interface{}) {

	if log.file != nil {

		fmt_log := fmt.Sprintf(f,params...)

		callTime := (float32)((float32)(conn.GetResponseTime().Sub(conn.GetRequestTime())) / (float32)(1000 * 1000))

		str := fmt.Sprintf("[%s] %s %s %d %.2f ms "  + "\n", conn.GetRequestTime().Format("2015-03-07 11:06:39 -0800 PST"),
			conn.GetRemoteAddr(), fmt_log,level, callTime)
		//log.file.Write([]byte(str))

		fmt.Printf("%s",str)
	}

}

