package main

import (
	"net"
	"log"
	"strings"

	pb "../../protocol"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"errors"
	"google.golang.org/grpc/credentials"
)

var grpcServer *grpc.Server

var bufferCutoff = 100

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

func (rpcserver *RPCServer) Attach(ctx context.Context, query *pb.ServerQuery) (*pb.ServerReply, error) {
	found := false
	for serv := range Servers {
		if serv == query.ProcessName {
			found = true
			break
		}
	}

	if !found {
		return &pb.ServerReply{}, errors.New("Process name not found.")
	}

	reply := &pb.ServerReply{}           //begin construction of reply
	server := Servers[query.ProcessName]

	//Parse ServerQuery object

	if query.MessageId == -1 { //client requests for latest messages
		reply.Messages = server.getLog(server.getLatestLogID()-bufferCutoff, server.getLatestLogID()+1)
		if server.getLatestLogID()-bufferCutoff >= 0 {
			reply.MessageId = uint64(server.getLatestLogID() - bufferCutoff)
		} else {
			reply.MessageId = 0
		}

	} else if query.MessageId > -1 { //client requests for specific message sets
		reply.Messages = server.getLog(int(query.MessageId-int64(bufferCutoff)), int(query.MessageId))
		if query.MessageId-int64(bufferCutoff) >= 0 {
			reply.MessageId = uint64(query.MessageId - int64(bufferCutoff))
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

	//send command to process
	if query.Command != "" {
		server.input(query.Command)
		server.addLog("Remote command executed: " + strings.Replace(query.Command, "\n", "", -1))
	}
	return reply, nil
}

func rpcserverStart() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 19005))
	if err != nil {
		addLog(err.Error())
		log.Fatal("Oh no! IPC listen error (check if the port has been taken):", err)
	}
	if instanceSettings.SSLEncryption {
		creds, err := credentials.NewServerTLSFromFile(instanceSettings.CertFilePath, instanceSettings.KeyFilePath)
		if err != nil {
			logFatal(err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(creds))
	} else {
		grpcServer = grpc.NewServer()
	}
	pb.RegisterRPCServerServer(grpcServer, &RPCServer{})
	grpcServer.Serve(lis)
}
