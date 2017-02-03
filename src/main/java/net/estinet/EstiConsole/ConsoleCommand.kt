package net.estinet.EstiConsole

import net.estinet.EstiConsole.commands.HelpCommand
import java.util.*

fun setupCommands(){
    commands.add(HelpCommand())
}

open class ConsoleCommand{
    // Name of Command
    var cName = "command"
    open fun run(args: ArrayList<String>){}
}
