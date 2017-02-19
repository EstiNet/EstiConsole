package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.EstiConsole
import java.util.*

class DebugCommand : ConsoleCommand() {
    init {
        super.cName = "debug"
        super.desc = "Toggles debug mode."
    }
    override fun run(args: ArrayList<String>){
        EstiConsole.debug = !EstiConsole.debug
    }
}