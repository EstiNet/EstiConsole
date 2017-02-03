package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.EstiConsole
import java.util.*

class StartCommand : ConsoleCommand() {
    init {
        super.cName = "start"
        super.desc = "Starts the java process. Turns on auto-start."
    }
    override fun run(args: ArrayList<String>){
        EstiConsole.autoStartOnStop = true
        if(EstiConsole.javaProcess!!.isAlive) EstiConsole.println("Java process already online!")
        else net.estinet.EstiConsole.startJavaProcess()
    }
}