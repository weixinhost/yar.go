package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/weixinhost/yar.go"
	"github.com/weixinhost/yar.go/packager"
)

type Server struct {
	class     interface{}
	methodMap map[string]string
	body      []byte
	Opt       *yar.Opt
	writer    io.Writer
}

func NewServer(class interface{}) *Server {
	server := new(Server)
	server.class = class
	server.methodMap = make(map[string]string, 32)
	server.Opt = yar.NewOpt()
	return server
}

func (server *Server) Register(rpcName string, methodName string) {
	server.log(yar.LogLevelDebug, "Register Handler %s %s", rpcName, methodName)
	server.methodMap[strings.ToLower(rpcName)] = methodName
}

func (server *Server) Handle(body []byte, writer io.Writer) *yar.Error {
	server.body = body
	server.writer = writer

	if len(server.body) < (yar.ProtocolLength + yar.PackagerLength) {
		return yar.NewError(yar.ErrorRequest, "request content errror:"+string(server.body))
	}

	header, err := server.readHeader()

	if err != nil {
		server.log(yar.LogLevelError, "[YarCall] readHeader error:%s", err.String())
		return err
	}

	request, err := server.readRequest(header)

	if err != nil {
		server.log(yar.LogLevelError, "[YarCall] readResponse error:%s", err.String())
		return err
	}

	response := yar.NewResponse()
	response.Status = yar.ERR_OKEY
	response.Protocol = header

	server.call(request, response)
	server.sendResponse(response)
	if response.Status != yar.ERR_OKEY {
		server.log(yar.LogLevelError, "[YarCall] %d %s Error:%s\n", request.Id, request.Method, response.Error)
		return yar.NewError(yar.ErrorResponse, response.Error)
	} else {
		server.log(yar.LoglevelNormal, "[YarCall] %d %s %s\n", request.Id, request.Method, "OKEY")
	}
	return nil
}

func (server *Server) readHeader() (*yar.Header, *yar.Error) {

	headerBuffer := bytes.NewBuffer(server.body[0 : yar.ProtocolLength+yar.PackagerLength])

	header := yar.NewHeaderWithBytes(headerBuffer)

	if header.MagicNumber != server.Opt.MagicNumber {
		return nil, yar.NewError(yar.ErrorProtocol, "magic number check failed.")
	}

	encrypt := server.Opt.Encrypt
	encryptKey := ""

	if header.Encrypt == 1 {
		if encrypt == false {
			return nil, yar.NewError(yar.ErrorProtocol, "this is a encrypt request,but server not support encrypt mode.")
		}

		if len(encryptKey) < 1 {
			return nil, yar.NewError(yar.ErrorProtocol, "this is a encrypt request,but server not set a encrypt private key")
		}
	}

	if header.Encrypt == 0 && encrypt == true {
		return nil, yar.NewError(yar.ErrorProtocol, "this server is encrypt,but request is not encrypt mode")
	}

	serviceName := string(header.Token[:])

	if server.Opt.CheckRequest {
		s1 := yar.StrToFixedBytes(server.Opt.ServiceName, 32)
		if serviceName != string(s1) {
			return nil, yar.NewError(yar.ErrorProtocol, "mismatch service name: "+string(s1)+"!="+serviceName)
		}
	}

	return header, nil
}

func (server *Server) readRequest(header *yar.Header) (*yar.Request, *yar.Error) {
	server.log(yar.LogLevelDebug, "[readRequest] %d %s %d %d", header.Id, header.Packager, header.MagicNumber, header.BodyLength)
	bodyLen := header.BodyLength
	bodyBuffer := server.body[90 : 90+bodyLen-8]

	request := yar.NewRequest()

	err := packager.Unpack(header.Packager[:], bodyBuffer, request)

	if err != nil {
		return nil, yar.NewError(yar.ErrorPackager, err.Error())
	}

	return request, nil
}

func (server *Server) sendResponse(response *yar.Response) *yar.Error {
	server.log(yar.LogLevelDebug, "[sendResponse] %d %d %s", response.Id, response.Status, fmt.Sprint(response.Retval))
	sendPackData, err := packager.Pack(response.Protocol.Packager[:], response)
	if err != nil {
		return yar.NewError(yar.ErrorResponse, err.Error())
	}
	response.Protocol.BodyLength = uint32(len(sendPackData) + 8)
	server.writer.Write(response.Protocol.Bytes().Bytes())
	server.writer.Write(sendPackData)
	return nil

}

func (server *Server) call(request *yar.Request, response *yar.Response) {

	defer func() {
		if r := recover(); r != nil {
			response.Status = yar.ERR_EMPTY_RESPONSE
			response.Error = "call handler internal panic:" + fmt.Sprint(r)
			if server.Opt.LogLevel&yar.LogLevelError > 0 {
				fmt.Println(r)
				debug.PrintStack()
			}
		}
	}()

	call_params := request.Params.([]interface{})

	class_fv := reflect.ValueOf(server.class)

	methodMap, ok := server.methodMap[strings.ToLower(request.Method)]

	var err bool

	if ok == false {
		_, err = class_fv.Type().MethodByName(request.Method)
		methodMap = request.Method
	} else {
		_, err = class_fv.Type().MethodByName(methodMap)
	}

	if err == false {
		response.Status = yar.ERR_EMPTY_RESPONSE
		response.Error = "call undefined api:" + request.Method
		return
	}

	fv := class_fv.MethodByName(methodMap)

	var real_params []reflect.Value

	if server.Opt.DynamicParam {
		real_params = make([]reflect.Value, fv.Type().NumIn())
	} else {

		if len(call_params) != fv.Type().NumIn() {
			response.Status = yar.ERR_EMPTY_RESPONSE
			response.Error = "mismatch handler param size"
			return
		}

		real_params = make([]reflect.Value, len(call_params))
	}

	func() {

		for i := 0; i < len(real_params); i++ {

			if i >= len(call_params) {
				tv := fv.Type().In(i).Kind()
				if tv == reflect.Ptr || tv == reflect.Map || tv == reflect.Array || reflect.Interface == tv {
					real_params[i] = reflect.New(fv.Type().In(i))
				} else {
					real_params[i] = reflect.Zero(fv.Type().In(i))
				}
				continue
			}

			tv := fv.Type().In(i).Kind()
			if tv == reflect.Ptr || tv == reflect.Map || tv == reflect.Array || reflect.Interface == tv {
				real_params[i] = reflect.New(fv.Type().In(i))
			}

			v := call_params[i]

			raw_val := reflect.ValueOf(v)

			if !raw_val.IsValid() {
				continue
			}

			//hack number
			if raw_val.Type().Name() == "Number" {

				fi := fv.Type().In(i)
				var coverErr error = nil
				verify := true
				nv := v.(json.Number)

				switch fi.Kind() {

				case reflect.Uint8:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(uint8(utv))
						break
					}

				case reflect.Uint16:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(uint16(utv))
						break
					}

				case reflect.Uint32:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(uint32(utv))
						break
					}

				case reflect.Uint64:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(uint64(utv))
						break
					}

				case reflect.Uint:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(uint(utv))
						break
					}

				case reflect.Int8:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(int8(utv))
						break
					}
				case reflect.Int16:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(int16(utv))
						break
					}
				case reflect.Int32:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(int32(utv))
						break
					}
				case reflect.Int64:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(int64(utv))
						break
					}
				case reflect.Int:
					{
						utv, err := nv.Int64()
						coverErr = err
						real_params[i] = reflect.ValueOf(int(utv))
						break
					}
				case reflect.Float32:
					{
						utv, err := nv.Float64()
						coverErr = err
						real_params[i] = reflect.ValueOf(float32(utv))
						break
					}
				case reflect.Float64:
					{
						utv, err := nv.Float64()
						coverErr = err
						real_params[i] = reflect.ValueOf(float64(utv))
						break
					}

				default:
					{
						verify = false
					}
				}

				if coverErr != nil {
					response.Status = yar.ERR_EMPTY_RESPONSE
					response.Error = "cover number type error:" + coverErr.Error()
					return
				}

				if verify == true {
					continue
				}

			}

			if raw_val.Type().Name() == "string" {

				var coverErr error = nil
				verify := true

				switch fv.Type().In(i).Kind() {

				case reflect.Uint8:
					{

						n, e := strconv.ParseUint(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(uint8(n))
						break
					}
				case reflect.Uint16:
					{

						n, e := strconv.ParseUint(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(uint16(n))
						break

					}
				case reflect.Uint32:
					{
						n, e := strconv.ParseUint(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(uint32(n))
						break

					}
				case reflect.Uint64:
					{

						n, e := strconv.ParseUint(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(uint64(n))
						break

					}
				case reflect.Uint:
					{

						n, e := strconv.ParseUint(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(uint(n))
						break

					}

				case reflect.Int8:
					{

						n, e := strconv.ParseInt(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(int8(n))
						break

					}
				case reflect.Int16:
					{

						n, e := strconv.ParseInt(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(int16(n))
						break

					}
				case reflect.Int32:
					{

						n, e := strconv.ParseInt(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(int32(n))
						break

					}
				case reflect.Int64:
					{

						n, e := strconv.ParseInt(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(int64(n))
						break

					}

				case reflect.Int:
					{

						n, e := strconv.ParseInt(raw_val.String(), 10, 64)
						coverErr = e
						real_params[i] = reflect.ValueOf(int(n))
						break

					}

				case reflect.Float32:
					{
						n, e := strconv.ParseFloat(raw_val.String(), fv.Type().In(i).Bits())
						coverErr = e
						real_params[i] = reflect.ValueOf(float32(n))
						break
					}

				case reflect.Float64:
					{
						n, e := strconv.ParseFloat(raw_val.String(), fv.Type().In(i).Bits())
						coverErr = e
						real_params[i] = reflect.ValueOf(float64(n))
						break
					}

				default:
					{
						verify = false
					}

				}

				if coverErr != nil {
					response.Status = yar.ERR_EMPTY_RESPONSE
					response.Error = "cover string to number error:" + coverErr.Error()
					return
				}

				if verify == true {
					continue
				}

			}

			real_params[i] = raw_val.Convert(fv.Type().In(i))
		}

		rs := fv.Call(real_params)
		if len(rs) < 1 {
			response.Return(nil)
			return
		}

		if len(rs) > 1 {
			response.Status = yar.ERR_EMPTY_RESPONSE
			response.Error = "unsupprted multi value return on rpc call"
			return
		}
		response.Return(rs[0].Interface())
	}()
}

func (server *Server) log(level int, logFmt string, v ...interface{}) {
	if level&server.Opt.LogLevel >= level {
		log.Printf(logFmt, v...)
		log.Println("")
	}
}
