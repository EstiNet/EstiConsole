package net.estinet.EstiConsole

class ShutdownHook : Runnable{
    override fun run() {
        EstiConsole.autoStartOnStop = false
        if(mode == Modes.BUNGEE) EstiConsole.sendJavaInput("end")
        else if (mode == Modes.SPIGOT || mode == Modes.PAPERSPIGOT) EstiConsole.sendJavaInput("stop")
        val thr: Thread = Thread({
            Thread.sleep(30000)
            if(EstiConsole.javaProcess.isAlive()) EstiConsole.javaProcess.destroyForcibly()
        })
        thr.start()
    }
}