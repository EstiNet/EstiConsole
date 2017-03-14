package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient

interface Message{
    val name: String;
    fun run(args: List<String>, session: SocketIOClient, ack: AckRequest)
}