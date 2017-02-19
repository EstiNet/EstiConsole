package net.estinet.EstiConsole.network

import io.scalecube.socketio.Session
import net.estinet.EstiConsole.password
import net.estinet.EstiConsole.sessionStorage
import net.estinet.EstiConsole.sessions

class HelloMessage : Message{
    override val name: String = "hello"
    override fun run(args: List<String>, session: Session) {
        if(args[0] == password){
            sessionStorage.put(session.sessionId, session)
            sessions.set(session.sessionId, true)
        }
    }
}