package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.SocketIOClient
import net.estinet.EstiConsole.EstiConsole
import java.io.File

class CurdirMessage : Message{
    override val name: String = "curdir"
    override fun run(args: List<String>, session: SocketIOClient) {
        var str = ""

        if(!File(args[0]).exists()){

        }
        session.sendEvent("curdir", EstiConsole.logByteArray)
    }
}
