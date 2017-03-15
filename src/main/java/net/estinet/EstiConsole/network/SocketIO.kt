package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.Configuration
import com.corundumstudio.socketio.SocketIOClient
import com.corundumstudio.socketio.SocketIOServer
import com.corundumstudio.socketio.listener.DataListener
import net.estinet.EstiConsole.*
import java.io.File
import java.io.FileInputStream
import java.security.KeyStore
import java.security.SecureRandom
import javax.net.ssl.KeyManagerFactory
import javax.net.ssl.SSLContext
import javax.net.ssl.TrustManagerFactory

object SocketIO {
    lateinit var sslServer: SocketIOServer
    fun doSocket() {
        net.estinet.EstiConsole.println("Starting server...")

        //EstiConsole.println("Initializing SSL context...")
        //val sslContext = initSslContext()

        val config: Configuration = Configuration();
        config.setPort(port)

        sslServer = SocketIOServer(config)
        networkOn = true
        sslServer.addConnectListener({client: SocketIOClient ->
            run {
                EstiConsole.println("Client has connected: " + client.remoteAddress)
                sessions.put(client.sessionId.toString(), false);
                client.sendEvent("event")
            }
        })
        for(message in messages){
            sslServer.addEventListener(message.name, String::class.java, {client: SocketIOClient, data, ack ->
                run{
                    try {
                        val str = data.toString()
                        if(EstiConsole.debug){
                            EstiConsole.println("Received: " + message.name + " " + str)
                        }
                        if (message.name == "hello" || sessions.get(client.sessionId.toString())!!) {
                            message.run(data.toString().split(" "), client, ack)
                        } else {
                            client.sendEvent("ecerror", "900")
                        }
                    }
                    catch(e: NullPointerException) {}
                }
            })
        }
        sslServer.addDisconnectListener { client: SocketIOClient ->
            run{
                EstiConsole.println("Client has disconnected: " + client.remoteAddress)
                sessions.remove(client.sessionId.toString())
                sessionStorage.remove(client.sessionId.toString())
            }
        }
        sslServer.start()
        net.estinet.EstiConsole.println("Started server on port $port!")
    }

    fun sendToAll(output: String) {
        for (s in sessionStorage.values) {
            s.sendEvent("log", output.toByteArray())
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