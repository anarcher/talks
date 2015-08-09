// +build OMIT

package main

import (
	"log"
	"net/http"
	"net/rpc"
	"time"
)

var (
	addr = "localhost:5001"
)

// S:HANDLER1 OMIT
type AddRequest struct {
	A, B int64
}
type AddResponse struct {
	V int64
}

type Handler struct{}

func (h Handler) Add(req AddRequest, res *AddResponse) error {
	res.V = req.A + req.B
	return nil
}

// E:HANDLER1 OMIT

// START1 OMIT

func main() {
	var add = new(Handler)
	s := rpc.NewServer()
	s.RegisterName("addsvc", add)
	s.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	go http.ListenAndServe(addr, s)
	time.Sleep(1 * time.Second)
	ret := ClientCall(1, 2)
	log.Println("ret:", ret)
}

// STOP1 OMIT

func ClientCall(a int64, b int64) int64 {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	res := new(AddResponse)
	req := AddRequest{a, b}
	err = client.Call("addsvc.Add", req, &res)
	if err != nil {
		log.Println("err", err)
		return 0
	}
	return res.V
}
