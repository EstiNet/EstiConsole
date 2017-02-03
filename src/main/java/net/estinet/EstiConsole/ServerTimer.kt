package net.estinet.EstiConsole

class ServerTimer : Runnable{
    override fun run() {
        Thread.sleep(3000)
        if(EstiConsole.autoStartOnStop && !EstiConsole.javaProcess!!.isAlive()){
            startJavaProcess()
        }
        run()
    }
}
