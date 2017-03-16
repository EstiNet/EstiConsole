package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient
import net.estinet.EstiConsole.EstiConsole
import java.io.File

class UploadGoodMessage : Message{
    override val name: String = "uploadgood"
    override fun run(args: List<String>, session: SocketIOClient, ack: AckRequest) {
        uploadInProgress.remove(session.remoteAddress.toString())
    }
}