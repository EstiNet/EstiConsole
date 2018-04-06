package main


func NetworkStart() {
	info("Starting client connection process...")
	go rpcserverStart()
}

