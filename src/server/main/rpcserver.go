package main

import (
	"net"
	"strings"

	pb "../../protocol"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"errors"
	"google.golang.org/grpc/credentials"
	"strconv"
)

var (
	grpcServer *grpc.Server

	bufferCutoff = 100

	invalidToken = errors.New("invalid authentication token")
)

/*
 * Define grpc struct for esticli method calls
 */

type RPCServer struct{}

func (rpcserver *RPCServer) Version(ctx context.Context, str *pb.StringRequest) (*pb.String, error) {
	if _, ok := checkToken(str.AuthToken); ok && instanceSettings.RequireAuth {
		return &pb.String{Str: version}, nil
	} else {
		return nil, invalidToken
	}
}

func (rpcserver *RPCServer) List(ctx context.Context, str *pb.StringRequest) (*pb.String, error) {
	if _, ok := checkToken(str.AuthToken); !ok && instanceSettings.RequireAuth { //TODO check permissions of user
		return nil, invalidToken
	}

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

func (rpcserver *RPCServer) Stop(ctx context.Context, str *pb.StringRequest) (*pb.String, error) {
	if _, ok := checkToken(str.AuthToken); !ok && instanceSettings.RequireAuth { //TODO check permissions of user
		return nil, invalidToken
	}

	output := StopClient(str.Str)
	if strings.Split(output, " ")[0] == "Stopped" {
		info(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) Start(ctx context.Context, str *pb.StringRequest) (*pb.String, error) {
	if _, ok := checkToken(str.AuthToken); !ok && instanceSettings.RequireAuth { //TODO check permissions of user
		return nil, invalidToken
	}

	output := StartClient(str.Str)
	if strings.Split(output, " ")[0] == "Started" {
		info(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) Kill(ctx context.Context, str *pb.StringRequest) (*pb.String, error) {
	if _, ok := checkToken(str.AuthToken); !ok && instanceSettings.RequireAuth { //TODO check permissions of user
		return nil, invalidToken
	}

	output := KillClient(str.Str)
	if strings.Split(output, " ")[0] == "Killed" {
		info(output)
	}
	return &pb.String{Str: output}, nil
}

func (rpcserver *RPCServer) InstanceStop(ctx context.Context, str *pb.StringRequest) (*pb.String, error) {
	if _, ok := checkToken(str.AuthToken); !ok && instanceSettings.RequireAuth { //TODO check permissions of user
		return nil, invalidToken
	}

	go Shutdown()
	return &pb.String{Str: "Host service shutting down."}, nil
}

func (rpcserver *RPCServer) Attach(ctx context.Context, query *pb.ServerQuery) (*pb.ServerReply, error) {
	if _, ok := checkToken(query.AuthToken); !ok && instanceSettings.RequireAuth { //TODO check permissions of user
		return nil, invalidToken
	}

	found := false
	for serv := range Servers {
		if serv == query.ProcessName {
			found = true
			break
		}
	}

	if !found {
		return &pb.ServerReply{}, errors.New("process name not found")
	}

	reply := &pb.ServerReply{} //begin construction of reply
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

func (rpcserver *RPCServer) Auth(ctx context.Context, query *pb.User) (*pb.String, error) {
	for _, user := range instanceSettings.Users {
		if user.Name == query.Name {
			if user.Password == query.Password {
				return &pb.String{Str: getNewToken(user.Name)}, nil
			} else {
				return nil, errors.New("invalid credentials")
			}
		}
	}
	return nil, errors.New("invalid credentials")
}

func rpcserverStart() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", instanceSettings.InstancePort))
	if err != nil {
		addLog(err.Error())
		logFatalStr("Oh no! IPC listen error (check if the port has been taken):" + err.Error())
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
	info("Starting RPCServer on port " + strconv.Itoa(int(instanceSettings.InstancePort)))
	grpcServer.Serve(lis)
}
