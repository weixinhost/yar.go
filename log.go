/**
 * A Simple Log System
 * We Need To Implement More Interface:
 * use bufio
 * write with async
 *
 */
package yar
import (
	"bytes"
	"os"
	"fmt"
	"time"
)

const (

	LOG_MAX_BUFFER_SIZE =  1024 * 1024 * 32 // default is 32MB

	LOG_BUFFER_FLUSH_INTERVAL_MS = 1000 * 5 // default is flush log every 5 seconds.

)

type Log struct
{

	buffers 	*bytes.Buffer

	path 		string

	fd 			*os.File

	buffer_len 	uint32

	last_flush	time.Time

}


func LogNew(path string) (*Log,error) {

	log := new(Log)

	log.path = path

	fd,err := os.OpenFile(log.path,os.O_CREATE | os.O_RDWR | os.O_APPEND,0775)

	if(err != nil){

		return nil,err

	}

	log.fd = fd

	log.buffers = new(bytes.Buffer)

	return log,err
}

func (self *Log) Normal(format string,l ...interface{}) bool {

	return self.writeBuffer("NORMAL",format,l...)

}

func (self *Log) Warning(format string,l ...interface{}) bool {

	return self.writeBuffer("WARNING",format,l...)

}

func (self *Log) Error(format string,l ...interface{}) bool {

	return self.writeBuffer("ERROR",format,l...)

}

func (self *Log) writeBuffer(level string, format string ,logs ...interface{}) bool {

	now_time := time.Now().Unix()

	time_str := time.Unix(now_time,0).Format("2006-01-02 15:04:05")

	ret := fmt.Sprintf("[%s] %s ",time_str,level)

	ret += fmt.Sprintf(format + "\n",logs...)

	self.buffers.WriteString(ret)

	//todo implement to async
	if(self.buffers.Len() >= LOG_MAX_BUFFER_SIZE ||
	   now_time - self.last_flush.Unix() > LOG_BUFFER_FLUSH_INTERVAL_MS){
		self.Flush()
	}

	return true

}

func (self *Log) Flush() bool {

	self.fd.Write(self.buffers.Bytes())
	self.buffers = new(bytes.Buffer)
	self.last_flush = time.Now()

	return true
}
