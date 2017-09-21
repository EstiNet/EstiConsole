package net.estinet.EstiConsole

class ShutdownHook : Runnable{
    override fun run() {
        println("Shutdown detected, ending wrapping process...")
        EstiConsole.autoStartOnStop = false
        if(mode == Modes.BUNGEE) EstiConsole.sendJavaInput("end")
        else if (mode == Modes.SPIGOT || mode == Modes.PAPERSPIGOT) EstiConsole.sendJavaInput("stop")
        val thr = Thread({
            Thread.sleep(30000)
            if(EstiConsole.javaProcess.isAlive()) EstiConsole.javaProcess.destroyForcibly()
        })
        thr.start()
        var i = 0
        while(i < 30){
            Thread.sleep(1000)
            if(!EstiConsole.javaProcess.isAlive()) return
            i++
        }
        if(EstiConsole.javaProcess.isAlive()) EstiConsole.javaProcess.destroyForcibly()
    }
}