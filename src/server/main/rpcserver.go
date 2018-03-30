package main

import (
	"net"
	"log"
	"strings"

	pb "../../protocol"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
)

var grpcServer *grpc.Server

/*
 * Define grpc struct for esticli method calls
 */

type RPCServer struct{}

func (rpcserver *RPCServer) Version(ctx context.Context, str *pb.String) (*pb.String, error) {
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
		info(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) Start(ctx context.Context, str *pb.String) (*pb.String, error) {
	output := StartClient(str.Str)
	if strings.Split(output, " ")[0] == "Started" {
		info(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) Kill(ctx context.Context, str *pb.String) (*pb.String, error) {
	output := KillClient(str.Str)
	if strings.Split(output, " ")[0] == "Killed" {
		info(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) InstanceStop(ctx context.Context, str *pb.String) (*pb.String, error) {
	go Shutdown()
	return &pb.String{Str: "Host service shutting down."}, nil
}

func (rpcserver *RPCServer) Attach(stream pb.RPCServer_AttachServer) error {
	for {
		query, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		reply := &pb.ServerReply{} //begin construction of reply
		server := Servers[query.ProcessName]

		//Parse ServerQuery object

		if query.MessageId == -1 { //client requests for latest messages
			reply.Messages = server.getLog(server.getLatestLogID()-100, server.getLatestLogID())
			if server.getLatestLogID()-100 >= 0 {
				reply.MessageId = uint64(server.getLatestLogID()-100)
			} else {
				reply.MessageId = 0
			}

		} else if query.MessageId > -1 { //client requests for specific message sets
			reply.Messages = server.getLog(int(query.MessageId-100), int(query.MessageId))
			if query.MessageId-100 >= 0 {
				reply.MessageId = uint64(query.MessageId-100)
			} else {
				reply.MessageId = 0
			}
		} else { //client doesn't require messages
			reply.MessageId = uint64(server.getLatestLogID())
			reply.Messages = []string{}
		}

		if query.GetCpu {
			reply.CpuUsage = GetCPUUsage()
		}
		if query.GetRam {
			reply.RamUsage = GetMemoryUsage() //TODO redo proc info
		}
		
		if err := stream.Send(reply); err != nil {
			info("[Error] " + err.Error())
			return err
		}
	}
}

func rpcserverStart() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 19005))
	if err != nil {
		addLog(err.Error())
		log.Fatal("Oh no! IPC listen error (check if the port has been taken):", err)
	}
	grpcServer = grpc.NewServer()
	pb.RegisterRPCServerServer(grpcServer, &RPCServer{})
	grpcServer.Serve(lis)
}
