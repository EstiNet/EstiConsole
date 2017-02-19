package net.estinet.EstiConsole.network

import io.netty.buffer.Unpooled
import io.scalecube.socketio.Session
import net.estinet.EstiConsole.EstiConsole

class CommandMessage : Message{
    override val name: String = "command"
    override fun run(args: List<String>, session: Session) {
        session.send(Unpooled.copiedBuffer(EstiConsole.logByteArray.toByteArray()))
    }
}