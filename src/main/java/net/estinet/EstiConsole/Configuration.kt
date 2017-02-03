package net.estinet.EstiConsole

import java.util.*
import java.io.*


fun setupConfiguration(){
    val f = File("esticonsole.properties")
    if(!f.exists())
        f.createNewFile()
    val prop = Properties()
    var output: OutputStream? = null
    val inputstream: InputStream? = null
    try {
        prop.load(inputstream)
        output = FileOutputStream("esticonsole.properties")
        if(prop.getProperty("port") == null) {
            println("Input EstiConsole port (not minecraft server port) (Default: 6921):")
            var input = System.console().readLine()
            if(input == "") input = "6921"
            prop.setProperty("port", input)
        }
        else if(prop.getProperty("server_jar_name") == null){
            println("Input server jar name (Default: minecraft_server.jar):")
            var input = System.console().readLine()
            if(input == "") input = "minecraft_server.jar"
            prop.setProperty("server_jar_name", input)
        }
        else if(prop.getProperty("mode") == null){
            println("Input the mode you want EstiConsole to run in (Default: SPIGOT):")
            var input = System.console().readLine()
            if(input == "") input = "SPIGOT"
            prop.setProperty("mode", input)
        }
        else if(prop.getProperty("password") == null){
            println("Input the password for EstiConsole connections (Default: pass123):")
            var input = System.console().readLine()
            if(input == "") input = "password"
            prop.setProperty("password", input)
        }
        prop.store(output, null)
    } catch (io: IOException) {
        io.printStackTrace()
    } finally {
        if (output != null) {
            try {
                output.close()
            } catch (e: IOException) {
                e.printStackTrace()
            }
            loadConfiguration()
        }
    }
}

fun loadConfiguration(){
    val prop = Properties()
    var input: InputStream? = null
    try {
        input = FileInputStream("config.properties")
        prop.load(input)
        port = prop.getProperty("port").toInt()
        serverJarName = prop.getProperty("server_jar_name")
        stmode = prop.getProperty("mode")
        password = prop.getProperty("password")
    } catch (ex: IOException) {
        ex.printStackTrace()
    } finally {
        if (input != null) {
            try {
                input.close()
            } catch (e: IOException) {
                e.printStackTrace()
            }

        }
    }
}
