package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.commands
import net.estinet.EstiConsole.disable
import java.util.*

class StopCommand() : ConsoleCommand() {
    init {
        super.cName = "stop"
        super.desc = "Stops the EstiConsole server."
    }
    override fun run(args: ArrayList<String>){
        disable()
    }
}

