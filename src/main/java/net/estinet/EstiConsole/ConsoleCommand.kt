package net.estinet.EstiConsole

import java.util.*

open class ConsoleCommand{
    // Name of Command
    var cName = "command"
    var desc = "I am a command"
    open fun run(args: ArrayList<String>){}
}
