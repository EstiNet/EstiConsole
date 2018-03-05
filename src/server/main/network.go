package main


func NetworkStart() {
	//server := tcp_server.New("localhost:9999")
	info("Starting esticli connection process...")
	go rpcserverStart()
	info("Started!")
	info("Starting client connection process...")
}