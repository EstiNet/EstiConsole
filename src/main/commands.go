package main

func CommandVersion(input string) {
	println("Version " + version)
}
func CommandInstanceStop(input string) {
	Shutdown()
}
func CommandHelp(input string) {
	println("-----Help-----")
	println("ec version          | Get the version of this instance.")
	println("ec instancestop     | Stop this instance of EstiConsole.")
	println("ec list             | List all of the client servers.")
	println("ec switch [process] | Switch view to another process.")
	println("ec stop [process]   | Stop the process using the default stop command.")
	println("ec start [process]  | Start the process.")
	println("ec kill [process]   | Forcibly kill the process.")
}
func CommandList(input string) {
	println("Clients:")
	for k, v := range Servers {
		var state string
		if v.IsOnline {
			state = "Online"
		} else {
			state = "Offline"
		}
		println(k + " (" + state + ")")
	}
}
func CommandSwitch(input string) {
	if _, ok := Servers[input]; ok {
		curServerView = Servers[input]
		ClearTerminal()
		print(Servers[input].Log)
		println("Successfully switched server view to " + input)
	} else {
		println("Server not found.")
	}
}
func CommandStop(input string) {
	if _, ok := Servers[input]; ok {
		Servers[input].AutoStart = false
		Servers[input].stop()
		println("Stopped " + Servers[input].Settings.InstanceName)
	} else {
		println("Server not found.")
	}
}
func CommandStart(input string) {
	if _, ok := Servers[input]; ok {
		if Servers[input].IsOnline {
			println("Process already online.")
		} else {
			Servers[input].AutoStart = true
			Servers[input].start()
			println("Started " + Servers[input].Settings.InstanceName)
		}
	} else {
		println("Server not found.")
	}
}
func CommandKill(input string) {
	if _, ok := Servers[input]; ok {
		if !Servers[input].IsOnline {
			println("Process is not online.")
		} else {
			Servers[input].AutoStart = false
			Servers[input].kill()
		}
	} else {
		println("Server not found.")
	}
}
