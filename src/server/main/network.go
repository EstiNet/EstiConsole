package main

import (
	pb "../../protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
	"strconv"
	"context"
	"google.golang.org/grpc/connectivity"
)

func NetworkStart() {
	info("Starting client connection process...")

	rpcserverStart()
	go func() { //Start connections with proxied processes and add them to map
		for i, server := range instanceSettings.ProxiedServers {
			cli, conn, token := StartRPCCon(&server)
			if server.Disabled {
				info("Skipping " + server.ProcessAlias + ".")
				continue
			}
			proxiedServerCon[server.ProcessAlias] = ProxiedServer{client: cli, connection: conn, token: token, config: server}
			instanceSettings.ProxiedServers[i] = server
		}
	}()
}


func StartRPCCon(server *ProxiedServerConfig) (client pb.RPCServerClient, conn *grpc.ClientConn, token string) {

	if server.Disabled {
		return
	}

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
				server.Disabled = true
				return
			}
		} else { //YAAAAAAAAAAAA encryption without mitm checks
			creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	info("Attempting proxy connection to " + server.ProcessAlias + "...")
	var err error
	conn, err = grpc.Dial(server.IP+":"+strconv.Itoa(int(server.Port)), opts...)
	if err != nil {
		info("Error connecting to process " + server.IP + ":" + strconv.Itoa(int(server.Port)) + ", is the address and port correct?:" + err.Error())
		return
	}
	client = pb.NewRPCServerClient(conn)

	tok, err := client.Auth(context.Background(), &pb.User{Name: server.Username, Password: server.Password})

	if err != nil && conn.GetState() == connectivity.TransientFailure { //if the state is failing ignore other errors
		info("Error connecting to process " + server.IP + ":" + strconv.Itoa(int(server.Port)) + ", is the address and port correct?:" + err.Error())
		return
	}

	if err != nil {
		info("Proxied process (" + server.ProcessAlias + ") authentication error: " + err.Error())
		server.Disabled = true
		return //prevent null dereference
	}
	info("Connected and authenticated with " + server.ProcessAlias + ".")
	token = tok.Str
	list, err := client.List(context.Background(), &pb.StringRequest{Str: "", AuthToken: token})
	if err != nil {
		info("Error: " + err.Error())
	} else {
		notFound := true
		for _, k := range list.Processes {
			if k.Name == server.ProcessName {
				notFound = false
				break
			}
		}
		if notFound {
			info("Proxied process (" + server.ProcessName + ") not found on server!")
			server.Disabled = true
		}
	}
	return
}
