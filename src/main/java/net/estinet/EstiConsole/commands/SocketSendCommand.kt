package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.EstiConsole
import net.estinet.EstiConsole.network.SocketIO
import java.util.*

class SocketSendCommand : ConsoleCommand() {
    init {
        super.cName = "socketsend"
        super.desc = "Sends a message to all connected websockets."
    }
    override fun run(args: ArrayList<String>){
        val arg = args.joinToString(" ")
        SocketIO.sendToAll(arg)
        EstiConsole.println("Sending $arg to all sockets.")
    }
}