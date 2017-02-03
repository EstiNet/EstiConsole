package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.EstiConsole
import net.estinet.EstiConsole.commands
import java.util.*

class StartCommand() : ConsoleCommand() {
    init {
        super.cName = "start"
        super.desc = "Starts the java process."
    }
    override fun run(args: ArrayList<String>){
        if(EstiConsole.javaProcess!!.isAlive) EstiConsole.println("Java process already online!")
        else net.estinet.EstiConsole.startJavaProcess()
    }
}