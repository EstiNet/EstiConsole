package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient
import java.io.File

class CurdirMessage : Message{
    override val name: String = "curdir"
    override fun run(args: List<String>, session: SocketIOClient, ack: AckRequest) {
        var str = ""

        if(!File(args[0]).exists()){
            ack.sendAckData("ecerror", "200")
        }
        else if(!File(args[0]).isDirectory){
            ack.sendAckData("ecerror", "202")
        }
        else {
            for(file in File(args[0]).listFiles()){
                if(file.isDirectory){
                    str += file.name + ":" + file.length()/1024/1024 + ":true"
                }
                else{
                    str += file.name + ":" + file.length()/1024/1024 + ":false"
                }
            }
            ack.sendAckData("curdir", str)
        }
    }
}
