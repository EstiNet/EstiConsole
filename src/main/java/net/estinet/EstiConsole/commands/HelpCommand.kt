package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.commands
import java.util.*

class HelpCommand : ConsoleCommand() {
    init {
        super.cName = "help"
        super.desc = "Displays help for EstiConsole commands."
    }
    override fun run(args: ArrayList<String>){
        println("----------EstiConsole Help----------")
        for(command in commands) println("${command.cName} : ${command.desc}")
    }
}
