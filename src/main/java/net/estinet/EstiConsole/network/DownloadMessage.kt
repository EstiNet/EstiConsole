package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient
import net.estinet.EstiConsole.EstiConsole
import java.io.File

class DownloadMessage : Message{
    override val name: String = "download"
    override fun run(args: List<String>, session: SocketIOClient, ack: AckRequest) {
        try {
            if (!File(args[0]).exists()) {
                if (EstiConsole.debug) {
                    EstiConsole.println("[Debug] download request: Doesn't exist. " + args[0])
                }
                ack.sendAckData("ecerror", "201")
            } else {
                if (EstiConsole.debug) {
                    EstiConsole.println("[Debug] download request: Sending " + args[0])
                }
                val f: File = File(args[0])
                val bytes = f.readBytes()
                if(bytes.size > 10000){
                    ack.sendAckData("download", bytes)
                }
                else {
                    ack.sendAckData("download", bytes)
                    ack.sendAckData("downloadgood")
                }
            }
        }
        catch(e: Throwable){
            e.printStackTrace()
            ack.sendAckData("ecerror", "901")
        }
    }
}