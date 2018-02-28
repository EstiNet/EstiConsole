package cli

import (
	"flag"
	"os"
)

var instanceName string
var args []string

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
	case "list":
		break
	case "status":
		break
	}
}
