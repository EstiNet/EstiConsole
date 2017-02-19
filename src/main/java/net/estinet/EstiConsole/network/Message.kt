package net.estinet.EstiConsole.network

import io.scalecube.socketio.Session
import java.util.*

interface Message{
    val name: String;
    fun run(args: List<String>, session: Session)
}