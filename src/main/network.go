package main


func NetworkStart() {
	//server := tcp_server.New("localhost:9999")
	println("Starting esticli connection process...")
	ipcserverStart()
	println("Started!")
	println("Starting client connection process...")
}