package main
import (
	"yar"
	"runtime"
	"flag"
	"time"
)

type snowflake_context struct  {

	init bool
	data_center_id int
	worker_id int
	sequence int64
	last_timestamp uint64
	twepoch uint64
	worker_id_bits uint
	data_center_id_bits uint
	sequence_bits uint
	worker_id_shift uint
	data_center_id_shift uint
	timestamp_shift uint
	sequence_mask int64
}


var context = snowflake_context{
	init : false,
}

func ukeyStartup(twepoch uint64,worker_id int,data_center_id int) (int){

	if worker_id == 0 || data_center_id == 0 {
		return -1;
	}

	if(context.init == true){
		return 0;
	}

	context.init = true
	context.twepoch = twepoch
	context.worker_id = worker_id
	context.data_center_id = data_center_id
	context.sequence = 0
	context.last_timestamp = 0
	context.worker_id_bits = 5
	context.data_center_id_bits = 5
	context.sequence_bits = 12

	context.worker_id_shift = uint(context.sequence_bits);
	context.data_center_id_shift = uint(context.sequence_bits + context.worker_id_bits);
	context.timestamp_shift = uint(context.sequence_bits + context.worker_id_bits + context.data_center_id_bits);
	context.sequence_mask = -1 ^ (-1 << context.sequence_bits);

	return 0;

}

func realTime() (uint64) {

	var retval uint64
	now := time.Now()
	retval = uint64((now.Unix() * 1000) + int64((now.Nanosecond() / 1000000)))
	return retval
}

func getUuid() (map[string]int64) {

	var retval map[string]int64 = make(map[string]int64,1)

	if ukeyStartup(1288834974657, *machine_id, *data_center) == -1 {

		retval["uuid"] = 0

		return retval
	}

	timestamp := realTime()

	if(context.last_timestamp == timestamp) {

		context.sequence = (context.sequence + 1) & context.sequence_mask;

		//溢出啦
		if(context.sequence == 0){

			time.Sleep(time.Nanosecond * 1000000)
		}

	}else{
		context.sequence = 0
	}

	context.last_timestamp = timestamp

	retval["uuid"] = int64(
		((uint64(timestamp) - uint64(context.twepoch)) << uint(context.timestamp_shift)) |
		(uint64(context.data_center_id) << uint64(context.data_center_id_shift)) |
		(uint64(context.worker_id) << uint64(context.worker_id_shift)) |
		uint64(context.sequence))

	return retval

}


var data_center = flag.Int("data-center",1,"the data center id")

var machine_id = flag.Int("worker_id",1,"the machine_id")

func main() {

	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	server,_:= yar.NewServer("http",":8088")
	server.RegisterHandler("uuid", getUuid)
	server.SetOpt(yar.SERVER_OPT_LOG_PATH,"/tmp/log.log")
	server.Serve()

}