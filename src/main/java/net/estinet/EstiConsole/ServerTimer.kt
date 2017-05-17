package net.estinet.EstiConsole

import net.estinet.EstiConsole.commands.RestartCommand
import java.util.*

class ServerTimer : Runnable{
    override fun run() {
        val thr = Thread({
            if(autoRestart == "yes") {
                while(true){
                    Thread.sleep(timeAutoRestart.toLong() * 60 * 60 * 1000)
                    RestartCommand().run(ArrayList())
                }
            }
        })
        thr.start()
        while(true){
            Thread.sleep(3000)
            if(EstiConsole.autoStartOnStop && !EstiConsole.javaProcess.isAlive()){
                startJavaProcess()
            }
        }
    }
}
