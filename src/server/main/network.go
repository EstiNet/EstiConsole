package main

import (
	pb "../../protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
	"strconv"
	"context"
)

func NetworkStart() {
	info("Starting client connection process...")

	rpcserverStart()
	go func() { //Start connections with proxied processes and add them to map
		for _, server := range instanceSettings.ProxiedServers {
			cli, conn, token := StartRPCCon(server)
			proxiedServerCon[server.ProcessName] = ProxiedServer{client: cli, connection: conn, token: token}
		}
	}()
}

//TODO check if the rpc connection went offline and auto repair
//TODO REGEN TOKEN AFTER 1 HOUR
//TODO CHECK IF PROCESS ACTUALLY EXISTS ON PROXIED

func StartRPCCon(server ProxiedServerConfig) (client pb.RPCServerClient, conn *grpc.ClientConn, token string) {
	var opts []grpc.DialOption

	if !server.HasTLS {
		opts = append(opts, grpc.WithInsecure()) //no encryption
	} else {
		// Create the client TLS credentials
		var creds credentials.TransportCredentials
		if server.CheckTLS { //encryption with IP SANs validation (for mitm attacks)
			var err error
			creds, err = credentials.NewClientTLSFromFile(server.CertFile, "")
			if err != nil {
				info("Could not load tls cert: " + err.Error())
			}
		} else { //YAAAAAAAAAAAA encryption without mitm checks
			creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	info("Attempting proxy connection to " + server.ProcessName + "...")
	var err error
	conn, err = grpc.Dial(server.IP+":"+strconv.Itoa(int(server.Port)), opts...)
	if err != nil {
		info("Error connecting to process " + server.IP + ":" + strconv.Itoa(int(server.Port)) + ", is the address and port correct?:" + err .Error())
	}
	client = pb.NewRPCServerClient(conn)
	tok, err := client.Auth(context.Background(), &pb.User{Name: server.Username, Password: server.Password})
	if err != nil {
		info("Proxied process (" + server.ProcessName + ") authentication error: " + err.Error())
	}
	token = tok.Str
	return
}
