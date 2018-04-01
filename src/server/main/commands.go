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
		for _, e := range Servers[input].Log {
			println(e)
		}
		println("Successfully switched server view to " + input)
	} else {
		println("Server not found.")
	}
}
func CommandStop(input string) {
	println(StopClient(input))
}
func CommandStart(input string) {
	println(StartClient(input))
}
func CommandKill(input string) {
	println(KillClient(input))
}
