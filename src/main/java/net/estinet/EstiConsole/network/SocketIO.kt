package net.estinet.EstiConsole.network

import io.netty.buffer.ByteBuf
import io.netty.buffer.Unpooled
import io.netty.util.CharsetUtil
import io.scalecube.socketio.Session
import io.scalecube.socketio.SocketIOListener
import io.scalecube.socketio.SocketIOServer
import java.io.FileInputStream
import java.security.KeyStore
import java.io.File
import java.security.SecureRandom
import javax.net.ssl.KeyManagerFactory
import javax.net.ssl.SSLContext
import javax.net.ssl.TrustManagerFactory
import io.scalecube.socketio.ServerConfiguration
import net.estinet.EstiConsole.*

object SocketIO {
    lateinit var sslServer: SocketIOServer
    fun doSocket() {
        net.estinet.EstiConsole.println("Starting server...")

        //EstiConsole.println("Initializing SSL context...")
        //val sslContext = initSslContext()

        val config = ServerConfiguration.builder()
                .port(port)
                .eventExecutorEnabled(false)
                .build()
        sslServer = SocketIOServer.newInstance(config)
        networkOn = true
        sslServer.listener = object : SocketIOListener {
            override fun onConnect(session: Session) {
                EstiConsole.println("Client has connected: " + session)
                sessions.put(session.sessionId, false);
            }

            override fun onMessage(session: Session, message: ByteBuf) {
                val str = message.toString(CharsetUtil.UTF_8)
                EstiConsole.println("Received: " + str)
                if (str.split(" ")[0] == "hello" || sessions.get(session.sessionId)!!) {
                    for (messaged in messages) {
                        if (messaged.name == str.split(" ")[0]) {
                            messaged.run(str.split(" ").subList(1, str.split(" ").size), session)
                        }
                    }
                    message.release()
                }
                else{
                    session.send(Unpooled.copiedBuffer("error 900".toByteArray()))
                }
            }

            override fun onDisconnect(session: Session) {
                EstiConsole.println("Client has disconnected: " + session)
                sessions.remove(session.sessionId)
                sessionStorage.remove(session.sessionId)
            }
        }
        sslServer.start()
        net.estinet.EstiConsole.println("Started server on port $port!")
    }

    fun sendToAll(output: String) {
        for (s in sessionStorage.values) {
            s.send(Unpooled.copiedBuffer(output.toByteArray()))
        }
    }

    @Throws(Exception::class)
    private fun initSslContext(): SSLContext {
        val keystorePassword = "password".toCharArray()
        val keystoreFile = File(SocketIO::class.java.getResource("/keystore.jks").getFile())
        val ks = KeyStore.getInstance("JKS")
        ks.load(FileInputStream(keystoreFile), keystorePassword)
        val kmf = KeyManagerFactory.getInstance("SunX509")
        kmf.init(ks, keystorePassword)
        val tmf = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm())
        tmf.init(ks)
        val sslContext = SSLContext.getInstance("TLS")
        sslContext.init(kmf.getKeyManagers(), tmf.getTrustManagers(), SecureRandom())
        return sslContext
    }
}