package net.estinet.EstiConsole

import net.estinet.EstiConsole.commands.HelpCommand
import java.util.*

fun setupCommands(){
    commands.add(HelpCommand())
}

open class ConsoleCommand{
    // Name of Command
    var cName = "command"
    var desc = "I am a command"
    open fun run(args: ArrayList<String>){}
}
