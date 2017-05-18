package net.estinet.EstiConsole

import java.io.IOException
import java.io.InputStreamReader
import java.io.BufferedReader
import java.io.InputStream

internal class LogStreamReader(`is`: InputStream, `stream`: InputStream) : Runnable {
    private val reader: BufferedReader
    private val reader2: BufferedReader
    init {
        this.reader = BufferedReader(InputStreamReader(`is`))
        this.reader2 = BufferedReader(InputStreamReader(`stream`))
    }
    override fun run() {
        val thr = Thread({otherRun()})
        thr.start()
        try {
            var line: String? = reader.readLine()
            while (line != null) {
                if(mode == Modes.BUNGEE)
                    line.replace("\n", "")
                println(line)
                parseJavaOutput(line)
                line = reader.readLine()
            }
            reader.close()
        } catch (e: IOException) {
            e.printStackTrace()
        }
    }
    fun otherRun(){
        try {
            var line: String? = reader2.readLine()
            while (line != null) {
                println("[ERROR] $line")
                line = reader2.readLine()
            }
            reader2.close()
        } catch (e: IOException) {
            e.printStackTrace()
        }
    }
}