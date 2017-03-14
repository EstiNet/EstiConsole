package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient
import io.scalecube.socketio.Session
import net.estinet.EstiConsole.processCommand

class CommandMessage : Message{
    override val name: String = "command"
    override fun run(args: List<String>, session: SocketIOClient, ack: AckRequest) {
        var str = ""
        for(s in args) str += s + " "
        processCommand(str)
    }
}