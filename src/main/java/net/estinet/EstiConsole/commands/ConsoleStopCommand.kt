package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.*
import java.util.*

class ConsoleStopCommand : ConsoleCommand() {
    init {
        super.cName = "consolestop"
        super.desc = "Stops the EstiConsole server."
    }
    override fun run(args: ArrayList<String>){
        if(mode == Modes.BUNGEE) EstiConsole.sendJavaInput("end")
        else if (mode == Modes.SPIGOT) EstiConsole.sendJavaInput("stop")
    }
}

