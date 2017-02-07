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

type HttpServer struct {
	handle map[string]*Server
}

func NewHttpServer() *HttpServer {
	server := new(HttpServer)
	server.handle = make(map[string]*Server, 0)
	return server
}

func (server *HttpServer) RegisterHandle(path string, h *Server) {
	server.handle[path] = h
}

func (server *HttpServer) Serve(addr string) {
	printlnWelcome()
	log.Println("Start Yar Server:" + addr)

	err := fasthttp.ListenAndServe(addr, server.innerHandle)

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

	body := ctx.PostBody()
	p := ctx.Path()
	path := string(p)

	h, ok := server.handle[path]

	if !ok {
		log.Println("No Yar Server Found On Path:" + path)
		return
	}

	buf := bytes.NewBufferString("")

	yarErr := h.Handle(body, buf)

	if yarErr != nil {
		log.Println("Yar Server Handle Error:" + yarErr.String())
	}

	bufBody := buf.Bytes()

	offset := 0

	for {
		n, err := ctx.Write(bufBody[offset:])
		if err != nil {
			log.Println("Yar Server Write Error:" + err.Error())
			break
		}
		offset += n
		if offset >= len(bufBody) {
			break
		}
	}
}
