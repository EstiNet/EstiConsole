package main

import "log"

func CommandHelp(input string) {
	println("-----Help-----")
	println("version          | Get the version of this instance.")
	println("status           | Get the status of the instance")
	println("instancestop     | Stop this instance of EstiConsole.")
	println("list             | List all of the client servers.")
	println("switch [process] | Switch view to another process.")
	println("stop [process]   | Stop the process using the default stop command.")
	println("start [process]  | Start the process.")
	println("kill [process]   | Forcibly kill the process.")
}

func CommandVersion(input string) {
	startCon()
	argss := Args{[]string{}}
	var reply string
	err := client.Call("Ipcserver.Version", argss, &reply)
	if err != nil {
		log.Fatal("ipcserver error:", err)
	}
	println("Version: ", reply)
}

func CommandList(input string) {
	startCon()
	argss := Args{[]string{}}
	var reply string
	err := client.Call("Ipcserver.List", argss, &reply)
	if err != nil {
		log.Fatal("ipcserver error:", err)
	}
	println(reply)
}

func CommandStop(input string) {
	startCon()
	argss := Args{[]string{input}}
	var reply string
	err := client.Call("Ipcserver.Stop", argss, &reply)
	if err != nil {
		log.Fatal("ipcserver error:", err)
	}
	println(reply)
}

func CommandInstanceStop(input string) {
	startCon()
	argss := Args{[]string{}}
	var reply string
	err := client.Call("Ipcserver.InstanceStop", argss, &reply)
	if err != nil {
		log.Fatal("ipcserver error:", err)
	}
	println(reply)
}

func CommandStatus(input string) {
	startCon()
	println("Connection successful!")
}