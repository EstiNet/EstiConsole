package main

import (
	"flag"
	"os"
	"log"
	"net/rpc"
)

type Args struct {
	Slice []string;
}

var instanceName string
var args []string

var client *rpc.Client

func startCon() *rpc.Client {
	println("Attempting connection to host server...")
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:19005")
	if err != nil {
		log.Fatal("Error connecting to host process, is it running?:", err)
	}
	return client
}

/*
 * Client entry point
 */
func main() {
	args = os.Args[1:]
	clientPtr := flag.String("instance", "s", "specify the instance to attach to")

	flag.Parse()

	instanceName = *clientPtr

	switch args[0] {
	case "help":
		println("-----Help-----")
		println("version          | Get the version of this instance.")
		println("status           | some sort of status thing")
		println("instancestop     | Stop this instance of EstiConsole.")
		println("list             | List all of the client servers.")
		println("switch [process] | Switch view to another process.")
		println("stop [process]   | Stop the process using the default stop command.")
		println("start [process]  | Start the process.")
		println("kill [process]   | Forcibly kill the process.")
		break
	case "version":
		client = startCon()
		argss := Args{[]string{}}
		var reply string
		err := client.Call("Ipcserver.Version", argss, &reply)
		if err != nil {
			log.Fatal("ipcserver error:", err)
		}
		println("Version: ", reply)
		break
	case "list":
		client = startCon()
		argss := Args{[]string{}}
		var reply string
		err := client.Call("Ipcserver.List", argss, &reply)
		if err != nil {
			log.Fatal("ipcserver error:", err)
			println(reply)
		}
		break
	case "status":
		break
	}

}