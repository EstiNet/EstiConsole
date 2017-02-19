package net.estinet.EstiConsole.network

import io.netty.buffer.Unpooled
import io.scalecube.socketio.Session
import net.estinet.EstiConsole.EstiConsole

class CurlogsMessage : Message{
    override val name: String = "curlogs"
    override fun run(args: List<String>, session: Session) {
        session.send(Unpooled.copiedBuffer(EstiConsole.logByteArray.toByteArray()))
    }
}
