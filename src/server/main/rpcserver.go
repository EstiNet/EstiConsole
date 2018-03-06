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

		if query.MessageId == -1 {
			reply.Messages = server.getLog(server.getLatestLogID()-100, server.getLatestLogID())
		} else if query.MessageId > -1 {
			reply.Messages = server.getLog(int(query.MessageId-100), int(query.MessageId))
		} else {
			reply.Messages = []string{}
		}

		//TODO get cpu and ram usage and process command

		if err := stream.Send(reply); err != nil {
			return err //TODO do something with error and parse
		}
	}
}

func rpcserverStart() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 19005))
	if err != nil {
		log.Fatal("Oh no! IPC listen error (check if the port has been taken):", err)
	}
	grpcServer = grpc.NewServer()
	pb.RegisterRPCServerServer(grpcServer, &RPCServer{})
	grpcServer.Serve(lis)
}
