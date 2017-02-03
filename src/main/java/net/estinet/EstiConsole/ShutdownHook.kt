package net.estinet.EstiConsole

class ShutdownHook : Runnable{
    override fun run() {
        if(mode == Modes.BUNGEE) EstiConsole.sendJavaInput("end")
        else if (mode == Modes.SPIGOT) EstiConsole.sendJavaInput("stop")
        var thr: Thread = Thread({
            Thread.sleep(30000)
            if(EstiConsole.javaProcess!!.isAlive()) EstiConsole.javaProcess!!.destroyForcibly()
        })
        thr.start()
    }
}