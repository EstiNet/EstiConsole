package net.estinet.EstiConsole.network

import io.scalecube.socketio.Session
import net.estinet.EstiConsole.processCommand

class CommandMessage : Message{
    override val name: String = "command"
    override fun run(args: List<String>, session: Session) {
        var str = ""
        for(s in args) str += s + " "
        processCommand(str)
    }
}