package net.estinet.EstiConsole

import java.util.*
import java.io.*

fun setupConfiguration() {
    val f = File("esticonsole.properties")
    if (!f.exists()) f.createNewFile()
    val prop = Properties()
    var output: FileOutputStream? = null
    var reader: InputStreamReader? = null
    try {
        reader = InputStreamReader(FileInputStream(f))
        prop.load(reader)
        output = FileOutputStream(f)
        if (prop.getProperty("port") == null) {
            println("Input EstiConsole port (not minecraft server port) (Default: 6921):")
            var input = System.console().readLine()
            if (input == "") input = "6921"
            prop.setProperty("port", input)
        }
        if (prop.getProperty("server_jar_name") == null) {
            println("Input server jar name (Default: minecraft_server.jar):")
            var input = System.console().readLine()
            if (input == "") input = "minecraft_server.jar"
            prop.setProperty("server_jar_name", input)
        }
        if (prop.getProperty("mode") == null) {
            println("Input the mode you want EstiConsole to run in (Default: SPIGOT):")
            var input = System.console().readLine()
            if (input == "") input = "SPIGOT"
            prop.setProperty("mode", input)
        }
        if (prop.getProperty("server_name") == null) {
            println("Input the server name (Default: Server):")
            var input = System.console().readLine()
            if (input == "") input = "Server"
            prop.setProperty("server_name", input)
        }
        if (prop.getProperty("password") == null) {
            println("Input the password for EstiConsole connections (Default: pass123):")
            var input = System.console().readLine()
            if (input == "") input = "pass123"
            prop.setProperty("password", input)
        }
        if (prop.getProperty("min_ram") == null) {
            println("Input the minimum amount of RAM for the server (Default: 512M):")
            var input = System.console().readLine()
            if (input == "") input = "512M"
            prop.setProperty("min_ram", input)
        }
        if (prop.getProperty("max_ram") == null) {
            println("Input the maximum amount of RAM for the server (Default: 2G):")
            var input = System.console().readLine()
            if (input == "") input = "2G"
            prop.setProperty("max_ram", input)
        }
    } catch (io: IOException) {
        io.printStackTrace()
    } finally {
        try {
            prop.store(output, null)
            if(reader != null) reader.close()
            if(output != null) output.close()
        } catch (e: IOException) {
            e.printStackTrace()
        }
        loadConfiguration()
    }
}

fun loadConfiguration() {
    val prop = Properties()
    var input: InputStream? = null
    try {
        input = FileInputStream("esticonsole.properties")
        prop.load(input)
        port = prop.getProperty("port").toInt()
        serverJarName = prop.getProperty("server_jar_name")
        stmode = prop.getProperty("mode")
        password = prop.getProperty("password")
        serverName = prop.getProperty("server_name")
        min_ram = prop.getProperty("min_ram")
        max_ram = prop.getProperty("max_ram")
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
