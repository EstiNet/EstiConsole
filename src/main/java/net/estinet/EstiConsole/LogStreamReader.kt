package net.estinet.EstiConsole

import java.io.IOException
import java.io.InputStreamReader
import java.io.BufferedReader
import java.io.InputStream

internal class LogStreamReader(`is`: InputStream) : Runnable {
    private val reader: BufferedReader
    init {
        this.reader = BufferedReader(InputStreamReader(`is`))
    }
    override fun run() {
        try {
            var line: String? = reader.readLine()
            while (line != null) {
                println(line)
                parseJavaOutput(line)
                line = reader.readLine()
            }
            reader.close()
        } catch (e: IOException) {
            e.printStackTrace()
        }
    }
}