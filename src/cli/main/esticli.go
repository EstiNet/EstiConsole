package main

import (
	"flag"
	_ "os"
	"log"
	"strings"
	"google.golang.org/grpc"

	pb "../../protocol"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
	"github.com/jroimartin/gocui"
	"context"
	"io/ioutil"
	"github.com/howeyc/gopass"
	"errors"
)

var (
	version = "v1.0.0"

	commands = make(map[string]interface{})

	args                                     []string
	address, port, certFile, user, masterKey *string
	verifyTLS                                *bool
	noTLS                                    *bool
	token                                    string

	password []byte

	conn *grpc.ClientConn

	client pb.RPCServerClient
)

/*
 * Initialize command functions
 */

func init() {
	commands["help"] = CommandHelp
	commands["version"] = CommandVersion
	commands["list"] = CommandList
	commands["status"] = CommandStatus
	commands["instancestop"] = CommandInstanceStop
	commands["stop"] = CommandStop
	commands["start"] = CommandStart
	commands["restart"] = CommandRestart
	commands["kill"] = CommandKill
	commands["attach"] = CommandAttach
}

/*
 * Client entry point
 */

func main() {

	//Initialize flags first
	getVer := flag.Bool("v", false, "get the version of the client")

	address = flag.String("ip", "127.0.0.1", "specify the address of the host")
	port = flag.String("p", "19005", "specify the port of the host")
	noTLS = flag.Bool("insecure", false, "specify whether or not to disable encryption")
	certFile = flag.String("cert", "none", "location of cert file (if using encryption)")
	verifyTLS = flag.Bool("verify", false, "whether or not to verify tls from server (if using encryption)")
	user = flag.String("username", "none", "specify the username to connect with (if using authentication)")
	masterKey = flag.String("masterkey", "none", "specify the location of the master key file (if using root authentication)")

	flag.Parse()       //Get the flag for user
	args = flag.Args() //os.Args[1:]

	if (*getVer) {
		println("EstiCli " + version)
	}

	//Check for command
	if len(args) == 0 {
		println("Unknown command, do /ec help.")
	} else if args[0] != "" {
		args[0] = strings.ToLower(args[0])
		found := false
		for k, v := range commands {
			if k == args[0] {
				in := ""
				for i, str := range args {
					if i != 0 {
						in += str
					}
				}
				v.(func(string))(in)
				found = true
				break
			}
		}
		if !found {
			println("Unknown command, do /ec help.")
		}
	}
}

/*
 * Handle RPC connection error
 */

func checkError(err error) {
	if err != nil {
		if cuiGUI != nil {
			(*cuiGUI).Update(func(g *gocui.Gui) error {
				return gocui.ErrQuit
			})
		}
		log.Fatal("[ERROR] ", err)
	}
}

func obtainCheckError(err error) bool {
	if err != nil && err == errors.New("invalid authentication token") {
		if *masterKey != "none" {
			dat, err := ioutil.ReadFile(*masterKey)
			if err != nil {
				log.Fatal(err)
			}
			tok, err := client.Auth(context.Background(), &pb.User{Name: "root", Password: string(dat)})
			if err != nil {
				log.Fatal(err)
			}
			token = tok.Str
		} else if *user != "none" {
			tok, err := client.Auth(context.Background(), &pb.User{Name: *user, Password: string(password)})
			if err != nil {
				log.Fatal(err)
			}
			token = tok.Str
		}
		return true
	} else {
		checkError(err)
	}
	return false
}

/*
 * Initialize connection with server
 */

func startCon() {
	var opts []grpc.DialOption

	if *noTLS {
		opts = append(opts, grpc.WithInsecure()) //no encryption
	} else {
		// Create the client TLS credentials
		var creds credentials.TransportCredentials
		if *verifyTLS { //encryption with IP SANs validation (for mitm attacks)
			var err error
			creds, err = credentials.NewClientTLSFromFile(*certFile, "")
			if err != nil {
				log.Fatal("Could not load tls cert: ", err)
			}
		} else { //YAAAAAAAAAAAA encryption without mitm checks
			creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	println("Attempting connection to host server...")
	var err error
	conn, err = grpc.Dial(*address + ":" + *port, opts...)
	if err != nil {
		log.Fatal("Error connecting to host process " + *address + ":" + *port+", is the address and port correct?:", err)
	}
	client = pb.NewRPCServerClient(conn)
	if *masterKey != "none" {
		dat, err := ioutil.ReadFile(*masterKey)
		if err != nil {
			log.Fatal(err)
		}
		tok, err := client.Auth(context.Background(), &pb.User{Name: "root", Password: string(dat)})
		if err != nil {
			log.Fatal(err)
		}
		token = tok.Str
	} else if *user != "none" {
		print("Password: ")
		var err error
		password, err = gopass.GetPasswd()
		if err != nil {
			log.Fatal(err)
		}
		tok, err := client.Auth(context.Background(), &pb.User{Name: *user, Password: string(password)})
		if err != nil {
			log.Fatal(err)
		}
		token = tok.Str
	}
}
