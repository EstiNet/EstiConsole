package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.EstiConsole
import net.estinet.EstiConsole.commands
import java.util.*

class KillCommand() : ConsoleCommand() {
    init {
        super.cName = "kill"
        super.desc = "Kills the java process (if online)."
    }
    override fun run(args: ArrayList<String>){
        println("Are you sure you want to kill the java process? (y/n)")
        val input = System.console().readLine()
        if(input.toLowerCase() == "y"){
            EstiConsole.javaProcess!!.destroy()
            EstiConsole.println("Killed the java process.")
        }
        else{
            EstiConsole.println("Prevented process kill.")
        }
    }
}

