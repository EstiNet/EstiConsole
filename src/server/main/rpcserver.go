package main

import (
	"strings"

	pb "../../protocol"
	"context"
	"google.golang.org/grpc"
	"errors"
	"net"
	"fmt"
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

func (rpcserver *RPCServer) List(ctx context.Context, str *pb.StringRequest) (*pb.ListResponse, error) {
	if _, ok := checkToken(str.AuthToken); !ok && instanceSettings.RequireAuth { //TODO check permissions of user
		return nil, invalidToken
	}

	ret := &pb.ListResponse{}
	for k, v := range Servers {
		var state string
		if v.IsOnline {
			state = "Online"
		} else {
			state = "Offline"
		}
		ret.Processes = append(ret.Processes, &pb.Process{Name: k, State: state})
	}

	for k, v := range proxiedServerCon {
		state := v.connection.GetState().String()
		if v.config.Disabled {
			state = "Disabled"
		}
		ret.Processes = append(ret.Processes, &pb.Process{Name: k, State: state})
	}

	return ret, nil
}

func (rpcserver *RPCServer) Stop(ctx context.Context, str *pb.StringRequest) (*pb.String, error) {
	if _, ok := checkToken(str.AuthToken); !ok && instanceSettings.RequireAuth { //TODO check permissions of user
		return nil, invalidToken
	}

	//check if to proxy the request to proxied process
	if sCon, ok := proxiedServerCon[str.Str]; ok {

		str, err := sCon.client.Stop(ctx, &pb.StringRequest{Str: sCon.config.ProcessName, AuthToken: sCon.token})
		if err != nil && err.Error() == "rpc error: code = Unknown desc = "+invalidToken.Error() { // regen the token if it's invalid
			RegenProxyToken(&sCon, str.Str)
			str, err = sCon.client.Stop(ctx, &pb.StringRequest{Str: sCon.config.ProcessName, AuthToken: proxiedServerCon[str.Str].token})
		}

		return str, err
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

	//check if to proxy the request to proxied process
	if sCon, ok := proxiedServerCon[str.Str]; ok {

		str, err := sCon.client.Start(ctx, &pb.StringRequest{Str: sCon.config.ProcessName, AuthToken: sCon.token})
		if err != nil && err.Error() == "rpc error: code = Unknown desc = "+invalidToken.Error() { // regen the token if it's invalid
			RegenProxyToken(&sCon, str.Str)
			str, err = sCon.client.Start(ctx, &pb.StringRequest{Str: sCon.config.ProcessName, AuthToken: proxiedServerCon[str.Str].token})
		}

		return str, err
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

	//check if to proxy the request to proxied process
	if sCon, ok := proxiedServerCon[str.Str]; ok {

		str, err := sCon.client.Kill(ctx, &pb.StringRequest{Str: sCon.config.ProcessName, AuthToken: sCon.token})
		if err != nil && err.Error() == "rpc error: code = Unknown desc = "+invalidToken.Error() { // regen the token if it's invalid
			RegenProxyToken(&sCon, str.Str)
			str, err = sCon.client.Kill(ctx, &pb.StringRequest{Str: sCon.config.ProcessName, AuthToken: proxiedServerCon[str.Str].token})
		}

		return str, err
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

	//check if to proxy the request to proxied process
	if sCon, ok := proxiedServerCon[query.ProcessName]; ok {
		query.AuthToken = sCon.token
		alias := query.ProcessName
		query.ProcessName = sCon.config.ProcessName

		reply, err := sCon.client.Attach(ctx, query)

		if err != nil && err.Error() == "rpc error: code = Unknown desc = "+invalidToken.Error() { // regen the token if it's invalid
			RegenProxyToken(&sCon, alias)
			query.AuthToken = proxiedServerCon[alias].token
			reply, err = sCon.client.Attach(ctx, query)
		}

		return reply, err
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

/*
 * Refresh a proxied server's token when it expires
 */

func RegenProxyToken(sCon *ProxiedServer, pName string) {
	tok, err := sCon.client.Auth(context.Background(), &pb.User{Name: sCon.config.Username, Password: sCon.config.Password})
	if err != nil {
		info("Proxied process (" + sCon.config.ProcessAlias + ") authentication error: " + err.Error())
	}
	info("Regenerated token with " + sCon.config.ProcessAlias + ".")
	ps := proxiedServerCon[pName]
	ps.token = tok.Str
	proxiedServerCon[pName] = ps
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
	go grpcServer.Serve(lis)
}
