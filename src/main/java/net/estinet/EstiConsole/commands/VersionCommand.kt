package net.estinet.EstiConsole.commands

import net.estinet.EstiConsole.ConsoleCommand
import net.estinet.EstiConsole.EstiConsole
import net.estinet.EstiConsole.commands
import java.util.*

class VersionCommand : ConsoleCommand(){
    init {
        super.cName = "version"
        super.desc = "Displays the version number for this instance of EstiConsole."
    }
    override fun run(args: ArrayList<String>){
        EstiConsole.println("EstiConsole ${EstiConsole.version}")
    }
}
