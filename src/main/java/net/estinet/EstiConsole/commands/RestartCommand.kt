package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.*
import java.util.*

class RestartCommand : ConsoleCommand() {
    init {
        super.cName = "restart"
        super.desc = "Restarts the Java process. Turns on auto-start."
    }
    override fun run(args: ArrayList<String>){
        EstiConsole.autoStartOnStop = true
        if(mode == Modes.BUNGEE) EstiConsole.sendJavaInput("end")
        else if (mode == Modes.SPIGOT || mode == Modes.PAPERSPIGOT) EstiConsole.sendJavaInput("stop")
    }
}
