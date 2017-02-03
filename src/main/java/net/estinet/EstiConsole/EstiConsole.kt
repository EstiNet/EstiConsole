package net.estinet.EstiConsole

import net.estinet.EstiConsole.commands.*
import java.util.*
import java.io.IOException
import java.io.File

object EstiConsole {
    var version: String = "v0.0.1-BETA"
    var javaProcess: Process? = null
    fun println(output: String) {
        System.out.println("${Locale.getLocale(LocaleType.PREFIX)} $output")
    }

    fun sendJavaInput(input: String) {
        val copyJavaProcess = javaProcess
        if (copyJavaProcess != null) copyJavaProcess.outputStream.bufferedWriter().write(input)
        else println("Oh noes! Can't send output to java process!")
    }
}

var mode: Modes = Modes.SPIGOT
var commands = ArrayList<ConsoleCommand>()

var port = 6921
var password = "pass123"
var serverJarName = "minecraft_server.jar"
var serverName = "Server"
var stmode = "SPIGOT"
var min_ram = "512M"
var max_ram = "2G"

/*
 * Command Initializer
 */
fun setupCommands() {
    commands.add(HelpCommand())
    commands.add(VersionCommand())
    commands.add(StopCommand())
}

/*
 * Program entry point.
 */

fun main(args: Array<String>) {
    System.out.println("EstiConsole.")
    Runtime.getRuntime().addShutdownHook(Thread(ShutdownHook()))
    System.out.println("Setting up Locale...")
    Locale.setupLocale()
    System.out.println(Locale.getLocale(LocaleType.ENABLING))
    enable()
}

fun enable() {
    /*
     * Startup Processes:
     */
    setupCommands()
    System.out.println("Setting up configuration...")
    setupConfiguration()
    var isMode = false
    for (value in Modes.values()) {
        if (stmode == value.toString()) {
            isMode = true
            mode = value
        }
    }
    if (isMode) {
        System.out.println("Mode selected: $mode")
        System.out.println("Welcome to EstiConsole.")
        val lambda = { startCommandProcess() }
        val thr: Thread = Thread(lambda)
        thr.start()
        EstiConsole.println("Starting Java process...")
        startJavaProcess()
    } else {
        println(Locale.getLocale(LocaleType.ERR_ON_START))
        EstiConsole.println("[Error] Incorrect mode specified!")
        EstiConsole.println("Exiting program...")
        System.exit(0)
    }
}

fun disable() {
    println(Locale.getLocale(LocaleType.DISABLING))
    /*
    * Disable Processes:
    */
    println(Locale.getLocale(LocaleType.DISABLED))
    System.exit(0)
}

fun startJavaProcess() {
    val pb = ProcessBuilder("java", "-Xms$min_ram", "-Xmx$max_ram", "-XX:+UseConcMarkSweepGC", "-XX:+UseParNewGC", "-XX:+CMSIncrementalPacing", "-XX:ParallelGCThreads=2", "-XX:+AggressiveOpts", "-d64", "-server", "-jar", serverJarName)
    pb.directory(File("./"))
    try {
        val p = pb.start()
        EstiConsole.javaProcess = p
        val lsr = LogStreamReader(p.inputStream)
        val thread = Thread(lsr, "LogStreamReader")
        thread.start()
    } catch (e: IOException) {
        e.printStackTrace()
    }
}

fun startCommandProcess() {
    while (true) {
        val input = System.console().readLine()
        System.out.println(input)
        val inputParsed = input.split(" ")
        if (inputParsed[0].toLowerCase() == "esticonsole" || inputParsed[0].toLowerCase() == "ec") {
            println("oh yea")
            var foundValue = false
            if (inputParsed.size < 2) {
                for (cc in commands) {
                    println(cc.cName)
                    if (cc.cName.toLowerCase() == inputParsed[1]) {
                        val args = ArrayList<String>()
                        val i = 0
                        while (i < inputParsed.size) {
                            if (i != 0 && i != 1) args.add(inputParsed[i])
                        }
                        cc.run(args)
                        foundValue = true
                        break
                    }
                }
            }
            if (!foundValue) println("Do /ec help for help!")
        } else {
            println("oh no")
            EstiConsole.sendJavaInput(input)
        }
    }
}

fun parseJavaOutput(output: String) {

}