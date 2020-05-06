package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

type httpRequest struct {
	method string
	path   string
}

type requestHandler struct {
	request httpRequest
}

func (h requestHandler) sleepAction() bool {
	return h.request.method == "GET" && h.request.path == "/sleep"
}

func (h requestHandler) connectionClose() bool {
	return h.request.method == "GET" && h.request.path == "/close"
}

func parseRequest(buf []byte) (httpRequest, error) {
	firstLine := strings.Split(string(buf), "\r\n")[0]
	splitedLine := strings.Split(firstLine, " ")

	if len(splitedLine) < 2 {
		return httpRequest{}, errors.New("can not parse the request")
	}

	request := httpRequest{
		method: splitedLine[0],
		path:   splitedLine[1],
	}

	return request, nil
}

func newRequestHandler(conn net.Conn) (requestHandler, error) {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return requestHandler{}, err
	}

	request, err := parseRequest(buf[:n])
	if err != nil {
		return requestHandler{}, err
	}

	return requestHandler{request: request}, nil
}

func execute(conn net.Conn) {
	defer conn.Close()

	handler, err := newRequestHandler(conn)
	if err != nil {
		fmt.Println("[Error] " + err.Error())
	}

	if handler.sleepAction() {
		fmt.Println("Sleeping...")
		time.Sleep(5 * time.Second)
		fmt.Println("Wake Up!!!")
	}

	if handler.connectionClose() {
		fmt.Println("Close")
		conn.Close()
		return
	}

	conn.Write([]byte("HTTP/1.1 200 OK\r\n"))
}

func listenAndServe(port string) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("[Error]" + err.Error())
	}

	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("[Error]" + err.Error())
		}

		go execute(conn)
	}
}

func main() {
	fmt.Println("Server starts at localhost:5454")
	listenAndServe(":5454")
}
