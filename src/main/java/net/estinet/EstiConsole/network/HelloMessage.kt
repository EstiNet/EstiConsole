package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.SocketIOClient
import io.netty.buffer.Unpooled
import net.estinet.EstiConsole.password
import net.estinet.EstiConsole.sessionStorage
import net.estinet.EstiConsole.sessions

class HelloMessage : Message{
    override val name: String = "hello"
    override fun run(args: List<String>, session: SocketIOClient) {
        if(args[0] == password){
            sessionStorage.put(session.sessionId.toString(), session)
            sessions.set(session.sessionId.toString(), true)
            session.sendEvent("authed")
        }
        else{
            session.sendEvent("error", "401")
        }
    }
}