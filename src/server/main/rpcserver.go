package main

import (
	"net"
	"log"
	"strings"

	pb "../../protocol"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type Args struct {
	Slice []string;
}

type RPCServer struct {

}

func (rpcserver *pb.RPCServerServer) Version(ctx context.Context, str *pb.String) (*pb.String, error) {
	return &pb.String{Str: version}, nil
}

func (rpcserver *RPCServer) List(ctx context.Context, str *pb.String) (*pb.String, error) {
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
	return &pb.String{Str: ret}, nil
}

func (rpcserver *RPCServer) Stop(ctx context.Context, str *pb.String) (*pb.String, error) {
	output := StopClient(str.Str)
	if strings.Split(output, " ")[0] == "Stopped" {
		println(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) Start(ctx context.Context, str *pb.String) (*pb.String, error) {
	output := StartClient(str.Str)
	println(strings.Split(output, " ")[0])
	if strings.Split(output, " ")[0] == "Started" {
		println(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) Kill(ctx context.Context, str *pb.String) (*pb.String, error) {
	output := KillClient(str.Str)
	if strings.Split(output, " ")[0] == "Killed" {
		println(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) InstanceStop(ctx context.Context, str *pb.String) (*pb.String, error) {
	go Shutdown()
	return &pb.String{Str: "Host service shutting down."}, nil
}

func rpcserverStart() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", "19005"))
	if err != nil {
		log.Fatal("Oh no! IPC listen error (check if the port has been taken):", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRPCServerServer(grpcServer, &RPCServer{})
	grpcServer.Serve(lis)
}
