package main
import ( "yar" )

func test_action(int_num int,float_num float32,str string)(string) {
	return "abcdefghj"
}

func main() {

	server := yar.NewServer("0.0.0.0", 6789)
	server.RegisterHandler("test", test_action)
	server.Run()

}
