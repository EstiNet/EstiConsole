package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.SocketIOClient
import net.estinet.EstiConsole.EstiConsole
import net.estinet.EstiConsole.password
import net.estinet.EstiConsole.sessionStorage
import net.estinet.EstiConsole.sessions

class HelloMessage : Message{
    override val name: String = "hello"
    override fun run(args: List<String>, session: SocketIOClient) {
        if(args[0] == password){
            sessionStorage.put(session.sessionId.toString(), session)
            sessions.set(session.sessionId.toString(), true)
            if(EstiConsole.debug){
                EstiConsole.println("${session.remoteAddress} has been authed.")
            }
            session.sendEvent("authed")
        }
        else{
            if(EstiConsole.debug){
                EstiConsole.println("${session.remoteAddress} has failed auth.")
            }
            session.sendEvent("error", "401")
        }
    }
}