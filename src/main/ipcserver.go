package main

import (
	"net/rpc"
	"net"
	"log"
	"net/http"
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
	if _, ok := Servers[arg.Slice[0]]; ok {
		Servers[arg.Slice[0]].AutoStart = false
		Servers[arg.Slice[0]].stop()
		output := "Stopped " + Servers[arg.Slice[0]].Settings.InstanceName
		println(output)
		*reply = output
	} else {
		*reply = "Server not found."
	}
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
