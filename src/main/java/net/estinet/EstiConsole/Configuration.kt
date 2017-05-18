package net.estinet.EstiConsole

import java.util.*
import java.io.*

fun setupConfiguration() {

    /*
     * Setup plugin files.
     */

    val update = File("update")
    if(!update.isDirectory()) update.mkdir()

    /*
     * Setup EstiConsole properties.
     */

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
            var input = console.readLine()
            if (input == "") input = "6921"
            prop.setProperty("port", input)
        }
        if (prop.getProperty("server_jar_name") == null) {
            println("Input server jar name (Default: minecraft_server.jar):")
            var input = console.readLine()
            if (input == "") input = "minecraft_server.jar"
            prop.setProperty("server_jar_name", input)
        }
        if (prop.getProperty("mode") == null) {
            println("Input the mode you want EstiConsole to run in (Default: SPIGOT):")
            var input = console.readLine()
            if (input == "") input = "SPIGOT"
            prop.setProperty("mode", input)
        }
        if (prop.getProperty("server_name") == null) {
            println("Input the server name (Default: Server):")
            var input = console.readLine()
            if (input == "") input = "Server"
            prop.setProperty("server_name", input)
        }
        if (prop.getProperty("password") == null) {
            println("Input the password for EstiConsole connections (Default: pass123):")
            var input = console.readLine()
            if (input == "") input = "pass123"
            prop.setProperty("password", input)
        }
        if (prop.getProperty("min_ram") == null) {
            println("Input the minimum amount of RAM for the server (Default: 512M):")
            var input = console.readLine()
            if (input == "") input = "512M"
            prop.setProperty("min_ram", input)
        }
        if (prop.getProperty("max_ram") == null) {
            println("Input the maximum amount of RAM for the server (Default: 2G):")
            var input = console.readLine()
            if (input == "") input = "2G"
            prop.setProperty("max_ram", input)
        }
        if (prop.getProperty("autoRestart") == null) {
            println("Input whether or not you want auto restart (Default: no):")
            var input = console.readLine()
            if (input == "") input = "no"
            prop.setProperty("autoRestart", input)
        }
        if (prop.getProperty("timeAutoRestart") == null) {
            println("Input time in hours for auto restart (Default: 24):")
            var input = console.readLine()
            if (input == "") input = "24"
            prop.setProperty("timeAutoRestart", input)
        }
        if (prop.getProperty("max_lines") == null) {
            fun verify(){
                println("Input max number of output lines to store in RAM (Default: 2000):")
                var input = console.readLine()
                try{
                    if(Integer.parseInt(input) < 1){
                        println("Please give a valid number.")
                        verify()
                    }
                    else{
                        if (input == "") input = "2000"
                        prop.setProperty("timeAutoRestart", input)
                    }
                }
                catch(e: Throwable){
                    net.estinet.EstiConsole.println("Please give a valid number.")
                    verify()
                }
            }
            verify()
        }
        if(prop.getProperty("amount_of_lines_to_cut_on_max") == null){
            fun verify(){
                println("Input number of lines to cut when line limit is reached (Default: 100):")
                var input = console.readLine()
                try{
                    if(Integer.parseInt(input) < 1){
                        println("Please give a valid number.")
                        verify()
                    }
                    else if(Integer.parseInt(input) <= Integer.parseInt("max_lines")){
                        println("The number must be smaller than the maximum amount of lines.")
                        verify()
                    }
                    else{
                        if (input == "") input = "100"
                        prop.setProperty("amount_of_lines_to_cut_on_max", input)
                    }
                }
                catch(e: Throwable){
                    net.estinet.EstiConsole.println("Please give a valid number.")
                    verify()
                }
            }
            verify()
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
        autoRestart = prop.getProperty("autoRestart")
        timeAutoRestart = prop.getProperty("timeAutoRestart")
        lineMax = prop.getProperty("max_lines")
        linesToCutOnMax = prop.getProperty("amount_of_lines_to_cut_on_max");
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
