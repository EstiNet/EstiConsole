package net.estinet.EstiConsole

class ShutdownHook : Runnable{
    override fun run() {
        if(mode == Modes.BUNGEE) EstiConsole.sendJavaInput("end")
        else if (mode == Modes.SPIGOT) EstiConsole.sendJavaInput("stop")
    }
}