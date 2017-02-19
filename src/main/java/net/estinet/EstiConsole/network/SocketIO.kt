package net.estinet.EstiConsole.network

import io.netty.buffer.ByteBuf
import io.netty.util.CharsetUtil
import io.scalecube.socketio.Session
import io.scalecube.socketio.SocketIOListener
import io.scalecube.socketio.SocketIOServer


object SocketIO{
    fun doSocket(){
        val logServer = SocketIOServer.newInstance(5000 /*port*/)
        logServer.listener = object : SocketIOListener {
            override fun onConnect(session: Session) {
                println("Connected: " + session)
            }

            override fun onMessage(session: Session, message: ByteBuf) {
                System.out.println("Received: " + message.toString(CharsetUtil.UTF_8))
                message.release()
            }

            override fun onDisconnect(session: Session) {
                println("Disconnected: " + session)
            }
        }
        logServer.start()
    }
}