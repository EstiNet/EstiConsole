package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient
import net.estinet.EstiConsole.EstiConsole
import java.io.File

class CurdirMessage : Message{
    override val name: String = "curdir"
    override fun run(args: List<String>, session: SocketIOClient, ack: AckRequest) {
        var str = ""

        if(!File(args[0]).exists()){
            ack.sendAckData("error", "200")
        }
        else if(!File(args[0]).isDirectory){
            ack.sendAckData("error", "202")
        }
        else {
            for(file in File(args[0]).listFiles()){
                str += file
            }
            session.sendEvent("curdir", EstiConsole.logByteArray)
        }
    }
}
