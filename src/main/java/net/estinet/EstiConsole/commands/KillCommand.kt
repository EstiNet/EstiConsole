package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.EstiConsole
import java.util.*

class KillCommand : ConsoleCommand() {
    init {
        super.cName = "kill"
        super.desc = "Kills the java process (if online). Turns off auto-start."
    }
    override fun run(args: ArrayList<String>){
        println("Are you sure you want to kill the java process? (y/n)")
        val input = System.console().readLine()
        if(input.toLowerCase() == "y"){
            EstiConsole.javaProcess!!.destroy()
            EstiConsole.autoStartOnStop = false
            EstiConsole.println("Killed the java process.")
        }
        else{
            EstiConsole.println("Prevented process kill.")
        }
    }
}

