package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.*
import java.util.*

class StopCommand : ConsoleCommand() {
    init {
        super.cName = "stop"
        super.desc = "Stops the java process (not the console). Turns off auto-start."
    }
    override fun run(args: ArrayList<String>){
        EstiConsole.autoStartOnStop = false
        if(mode == Modes.BUNGEE) EstiConsole.sendJavaInput("end")
        else if (mode == Modes.SPIGOT || mode == Modes.PAPERSPIGOT) EstiConsole.sendJavaInput("stop")
    }
}

