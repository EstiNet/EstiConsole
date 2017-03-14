package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient
import net.estinet.EstiConsole.EstiConsole

class CurlogsMessage : Message{
    override val name: String = "curlogs"
    override fun run(args: List<String>, session: SocketIOClient, ack: AckRequest) {
        session.sendEvent("curlogs", EstiConsole.logByteArray)
    }
}
