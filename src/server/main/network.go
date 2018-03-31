package main


func NetworkStart() {
	info("Starting esticli connection process...")
	go rpcserverStart()
	info("Started!")
	info("Starting client connection process...")
}