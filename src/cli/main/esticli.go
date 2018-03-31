package main

import (
	"flag"
	_ "os"
	"log"
	"strings"
	"google.golang.org/grpc"

	pb "../../protocol"
)

var version = "v1.0.0"

var commands = make(map[string]interface{})

var args []string
var address *string
var port *string

var conn *grpc.ClientConn

var client pb.RPCServerClient

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
	commands["kill"] = CommandKill
	commands["attach"] = CommandAttach
}

/*
 * Client entry point
 */

func main() {

	//Initialize flags first
	getVer := flag.Bool("v", false, "get the version of the client")

	address = flag.String("a", "127.0.0.1", "specify the address of the host")
	port = flag.String("p", "19005", "specify the port of the host")

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
		log.Fatal("[ERROR] ", err)
	}
}

/*
 * Initialize connection with server
 */

func startCon() {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	println("Attempting connection to host server...")
	var err error
	conn, err = grpc.Dial(*address + ":" + *port, opts...)
	if err != nil {
		log.Fatal("Error connecting to host process " + *address + ":" + *port+", is the address and port correct?:", err)
	}
	client = pb.NewRPCServerClient(conn)
}
