package main

import (
	"net/rpc"
	"net"
	"log"
	"net/http"
	"strings"
)

type Args struct {
	Slice []string;
}

type Ipcserver string

func (ipcserver *Ipcserver) Version(arg *Args, reply *string) error {
	*reply = version
	return nil
}

func (ipcserver *Ipcserver) List(arg *Args, reply *string) error {
	ret := ""
	ret += "Clients:\n"
	for k, v := range Servers {
		var state string
		if v.IsOnline {
			state = "Online"
		} else {
			state = "Offline"
		}
		ret += k + " (" + state + ")\n"
	}
	*reply = ret
	return nil
}

func (ipcserver *Ipcserver) Stop(arg *Args, reply *string) error {
	output := StopClient(arg.Slice[0])
	if strings.Split(output, " ")[0] == "Stopped" {
		println(output)
	}
	*reply = output
	return nil
}

func (ipcserver *Ipcserver) Start(arg *Args, reply *string) error {
	output := StartClient(arg.Slice[0])
	if strings.Split(output, " ")[0] == "Started" {
		println(output)
	}
	*reply = output
	return nil
}

func (ipcserver *Ipcserver) Kill(arg *Args, reply *string) error {
	output := KillClient(arg.Slice[0])
	if strings.Split(output, " ")[0] == "Killed" {
		println(output)
	}
	*reply = output
	return nil
}

func (ipcserver *Ipcserver) InstanceStop(arg *Args, reply *string) error {
	go Shutdown()
	*reply = "Host service shutting down."
	return nil
}

func ipcserverStart() {
	ipcserver := new(Ipcserver)
	rpc.Register(ipcserver)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":19005")
	if e != nil {
		log.Fatal("Oh no! IPC listen error (check if the port has been taken):", e)
	}
	go http.Serve(l, nil)
}
