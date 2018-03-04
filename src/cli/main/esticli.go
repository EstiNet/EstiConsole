package main

import (
	"flag"
	"os"
	"log"
	"net/rpc"
	"strings"
)

type Args struct {
	Slice []string;
}

var version = "v1.0.0"

var commands = make(map[string]interface{})

var args []string
var address *string
var port *string

var client *rpc.Client

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
}

/*
 * Client entry point
 */

func main() {
	args = os.Args[1:]

	//Initialize flags first
	getVer := flag.Bool("v", false, "get the version of the client")

	address = flag.String("a", "127.0.0.1", "specify the address of the host")
	port = flag.String("p", "19005", "specify the port of the host")

	flag.Parse() //Get the flag for user

	if (*getVer) {
		println("EstiCli " + version)
	}

	args[0] = strings.ToLower(args[0])

	//Check for command
	if args[0] != "" {
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
 * Handle IPC connection error
 */

func checkError(err error) {
	if err != nil {
		log.Fatal("IPC Error", err)
	}
}

/*
 * Initialize connection with server
 */

func startCon() {
	println("Attempting connection to host server...")
	clienti, err := rpc.DialHTTP("tcp", *address + ":" + *port)
	if err != nil {
		log.Fatal("Error connecting to host process " + *address + ":" + *port+", is the address and port correct?:", err)
	}
	client = clienti
}
