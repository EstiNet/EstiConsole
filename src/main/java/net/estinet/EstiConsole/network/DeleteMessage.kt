package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient
import net.estinet.EstiConsole.EstiConsole
import java.io.File

class DeleteMessage : Message{
    override val name: String = "delete"
    override fun run(args: List<String>, session: SocketIOClient, ack: AckRequest) {
        try {
            if (!File(args[0]).exists()) {
                if (EstiConsole.debug) {
                    EstiConsole.println("[Debug] delete request: Doesn't exist. " + args[0])
                }
                ack.sendAckData("ecerror", "201")
            } else {
                if (EstiConsole.debug) {
                    EstiConsole.println("[Debug] delete request: Deleted " + args[0])
                }
                File(args[0]).delete()
                ack.sendAckData("good")
            }
        }
        catch(e: Throwable){
            e.printStackTrace()
            ack.sendAckData("ecerror", "901")
        }
    }
}