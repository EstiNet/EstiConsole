package main

import _ "log"

func CommandHelp(input string) {
	println("-----Help-----")
	println("version          | Get the version of this instance.")
	println("status           | Get the status of the instance")
	println("instancestop     | Stop this instance of EstiConsole.")
	println("list             | List all of the client servers.")
	println("attach [process] | Switch view to another process.")
	println("stop [process]   | Stop the process using the default stop command.")
	println("start [process]  | Start the process.")
	println("kill [process]   | Forcibly kill the process.")
}

func CommandVersion(input string) {
	startCon()
	argss := Args{[]string{}}
	var reply string
	err := client.Call("Ipcserver.Version", argss, &reply)
	checkError(err)
	println("Version: ", reply)
}

func CommandList(input string) {
	startCon()
	argss := Args{[]string{}}
	var reply string
	err := client.Call("Ipcserver.List", argss, &reply)
	checkError(err)
	println(reply)
}

func CommandStop(input string) {
	startCon()
	argss := Args{[]string{input}}
	var reply string
	err := client.Call("Ipcserver.Stop", argss, &reply)
	checkError(err)
	println(reply)
}

func CommandInstanceStop(input string) {
	startCon()
	argss := Args{[]string{}}
	var reply string
	err := client.Call("Ipcserver.InstanceStop", argss, &reply)
	checkError(err)
	println(reply)
}

func CommandStart(input string) {
	startCon()
	argss := Args{[]string{input}}
	var reply string
	err := client.Call("Ipcserver.Start", argss, &reply)
	checkError(err)
	println(reply)
}

func CommandKill(input string) {
	startCon()
	argss := Args{[]string{input}}
	var reply string
	err := client.Call("Ipcserver.Kill", argss, &reply)
	checkError(err)
	println(reply)
}

func CommandAttach(input string) {

}

func CommandStatus(input string) {
	startCon()
	println("Connection successful!")
}
