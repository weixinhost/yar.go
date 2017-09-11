package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime/debug"

	"github.com/valyala/fasthttp"
)

func printlnWelcome() {

	name := `
 __   __   ___   ______ 
\ \ / /  / _ \  | ___ \
 \ V /  / /_\ \ | |_/ /
  \ /   |  _  | |    / 
  | |   | | | | | |\ \ 
  \_/   \_| |_/ \_| \_|

`
	welcome := "Welcome to Yar Server...\n" + name + "\n"
	fmt.Println(welcome)
}

type Handle func(body []byte, writer io.Writer)

type DirectHandle func(*fasthttp.RequestCtx, []byte, io.Writer)
type InitServer func() *Server

type HttpServer struct {
	handle       map[string]InitServer
	directHandle map[string]DirectHandle
}

func NewHttpServer() *HttpServer {
	server := new(HttpServer)
	server.handle = make(map[string]InitServer, 0)
	server.directHandle = make(map[string]DirectHandle, 0)
	return server
}

func (server *HttpServer) RegisterHandle(path string, h InitServer) {
	server.handle[path] = h
}

func (server *HttpServer) RegisterDirectHandle(path string, h DirectHandle) {
	server.directHandle[path] = h
}

func (server *HttpServer) Serve(addr string) {
	printlnWelcome()
	log.Println("Start Yar Server:" + addr)
	h := fasthttp.CompressHandler(server.innerHandle)
	err := fasthttp.ListenAndServe(addr, h)

	if err != nil {
		log.Println("Start Yar Server Error:" + err.Error())
		return
	}
}

func (server *HttpServer) innerHandle(ctx *fasthttp.RequestCtx) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			debug.PrintStack()
		}
	}()

	var body []byte

	if (string)(ctx.Request.Header.Peek("Content-Encoding")) == "gzip" {
		temp, err := ctx.Request.BodyGunzip()
		if err != nil {
			log.Println("Parse Body In Gunzip failed:" + err.Error())
			return
		}
		body = temp
	} else {
		body = ctx.PostBody()
	}
	p := ctx.Path()
	path := string(p)

	buf := bytes.NewBufferString("")

	if h, ok := server.directHandle[path]; ok {
		h(ctx, body, buf)
	} else {
		hf, ok := server.handle[path]

		if !ok {
			log.Println("No Yar Server Found On Path:" + path)
			return
		}
		h := hf()
		yarErr := h.Handle(body, buf)

		if yarErr != nil {
			log.Println("Yar Server Handle Error:" + yarErr.String())
		}
	}
	bufBody := buf.Bytes()
	var err error
	current := 0
	for {
		n, err := ctx.Write(bufBody[current:])
		if err != nil {
			break
		}
		current += n
		if current >= len(bufBody) {
			break
		}
	}

	if err != nil {
		log.Println("Yar Server Send Error:" + err.Error())
	}
}
