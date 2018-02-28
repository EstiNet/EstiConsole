package main

import (
	"net/rpc"
	"net"
	"log"
	"net/http"
	"github.com/phayes/freeport"
)

type Ipcserver string

func (ipcserver *Ipcserver) Version(arg *string, reply *string) error {
	*reply = version
	return nil
}

func ipcserverStart() {
	ipcserver := new(Ipcserver)
	rpc.Register(ipcserver)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":7269")
	if e != nil {
		log.Fatal("Oh no! IPC listen error (check if the port has been taken):", e)
	}
	go http.Serve(l, nil)
}
