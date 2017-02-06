package net.estinet.EstiConsole

import jline.console.ConsoleReader
import jline.console.CursorBuffer
import net.estinet.EstiConsole.commands.*
import org.fusesource.jansi.AnsiConsole
import java.io.File
import java.io.IOException
import java.io.OutputStreamWriter
import java.io.PrintWriter
import java.net.InetSocketAddress
import java.nio.channels.SelectionKey
import java.nio.channels.Selector
import java.nio.channels.ServerSocketChannel
import java.nio.channels.SocketChannel
import java.nio.file.Files
import java.util.*


object EstiConsole {
    var version: String = "v1.0.0"
    var javaProcess: Process? = null
    var writer: PrintWriter? = null
    var autoStartOnStop = false

    fun println(output: String) {
        stashLine()
        System.out.println("${Locale.getLocale(LocaleType.PREFIX)} $output")
        unstashLine()
    }
    fun sendJavaInput(input: String) {
        try {
            writer?.write("$input\n")
            writer?.flush()
            writer = PrintWriter(OutputStreamWriter(javaProcess?.outputStream))
        } catch(e: NullPointerException){
            println("Oh noes! Can't send output to java process! Is it offline?")
        }
    }
}

var mode: Modes = Modes.SPIGOT
var commands = ArrayList<ConsoleCommand>()

var channel: ServerSocketChannel? = null

var port = 6921
var password = "pass123"
var serverJarName = "minecraft_server.jar"
var serverName = "Server"
var stmode = "SPIGOT"
var min_ram = "512M"
var max_ram = "2G"

var console: ConsoleReader = ConsoleReader()
private var stashed: CursorBuffer? = null

/*
 * Command Initializer
 */
fun setupCommands() {
    commands.add(HelpCommand())
    commands.add(VersionCommand())
    commands.add(StopCommand())
    commands.add(StartCommand())
    commands.add(KillCommand())
    commands.add(ConsoleStopCommand())
    commands.add(RestartCommand())
}

/*
 * Program entry point.
 */

fun main(args: Array<String>) {
    println("EstiConsole.")
    Runtime.getRuntime().addShutdownHook(Thread(ShutdownHook()))
    AnsiConsole.systemInstall()
    println("Setting up Locale...")
    Locale.setupLocale()
    println(Locale.getLocale(LocaleType.ENABLING))
    enable()
}

fun enable() {
    /*
     * Startup Processes:
     */
    setupCommands()
    println("Setting up configuration...")
    setupConfiguration()
    var isMode = false
    for (value in Modes.values()) {
        if (stmode == value.toString()) {
            isMode = true
            mode = value
        }
    }
    if (isMode) {
        println("Mode selected: $mode")
        println("Welcome to EstiConsole.")
        val lambda = { startCommandProcess() }
        val thr: Thread = Thread(lambda)
        thr.start()
        EstiConsole.println("Starting Java process...")
        startJavaProcess()
        val th1 = Thread(ServerTimer())
        th1.start()
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
    EstiConsole.autoStartOnStop = false
    AnsiConsole.systemUninstall()
    println(Locale.getLocale(LocaleType.DISABLED))
    System.exit(0)
}

private fun startJavaProcessPluginFetch(){
    val update = File("./update")
    for(file in update.listFiles()){
        val p = File("plugins/${file.name}")
        //println("${file.name} ${File("plugins/${file.name}")} $serverJarName")
        if(serverJarName == file.name){
            EstiConsole.println("Updating java process jar $serverJarName...")
            File(serverJarName).delete()
            Files.copy(file.toPath(), File(file.name).toPath())
            file.delete()
        }
        else if(p.exists()){
            EstiConsole.println("Updating plugin $file...")
            p.delete()
            Files.copy(file.toPath(), p.toPath())
            file.delete()
        }
        else{
            Files.copy(file.toPath(), p.toPath())
            file.delete()
        }
    }
}

fun startJavaProcess() {
    EstiConsole.println("Fetching update folder...")
    startJavaProcessPluginFetch()
    EstiConsole.println("Starting jar...")
    val pb = ProcessBuilder("java", "-Xms$min_ram", "-Xmx$max_ram", "-XX:+UseConcMarkSweepGC", "-XX:+UseParNewGC", "-XX:+CMSIncrementalPacing", "-XX:ParallelGCThreads=2", "-XX:+AggressiveOpts", "-d64", "-server", "-jar", serverJarName, "-o true")
    pb.directory(File("./"))
    try {
        val process: Process = pb.start()
        EstiConsole.javaProcess = process
        EstiConsole.writer = PrintWriter(OutputStreamWriter(EstiConsole.javaProcess?.outputStream))
        val lsr = LogStreamReader(process.inputStream, process.errorStream)
        val thread = Thread(lsr, "LogStreamReader")
        thread.start()
        EstiConsole.autoStartOnStop = true
    } catch (e: IOException) {
        e.printStackTrace()
    }
}

fun startCommandProcess() {
    while (true) {
        console.setPrompt(">");
        val input = console.readLine()
        val inputParsed = input.split(" ")
        if (inputParsed[0].toLowerCase() == "esticonsole" || inputParsed[0].toLowerCase() == "ec") {
            var foundValue = false
            if (inputParsed.size >= 2) {
                for (cc in commands) {
                    if (cc.cName.toLowerCase() == inputParsed[1]) {
                        val args = ArrayList<String>()
                        var i = 0
                        while (i < inputParsed.size) {
                            if (i != 0 && i != 1) args.add(inputParsed[i])
                            i++
                        }
                        cc.run(args)
                        foundValue = true
                        break
                    }
                }
            }
            if (!foundValue) println("Do /ec help for help!")
        } else {
            EstiConsole.sendJavaInput(input)
        }
    }
}

fun parseJavaOutput(output: String) {

}

fun startNetworkProcess() {
    val selector = Selector.open()

    channel = ServerSocketChannel.open()
    val address = InetSocketAddress("localhost", port)
    channel.apply {
        channel!!.bind(address)
        channel!!.configureBlocking(false)
        channel!!.register(selector, channel!!.validOps())
    }

    while (true) {
        selector.select()

        val keys = selector.selectedKeys()
        val iterator = keys.iterator()
        while (iterator.hasNext()) {
            val key = iterator.next()
            if (key.isAcceptable) {
                channel!!.accept().apply {
                    configureBlocking(false)
                    register(selector, SelectionKey.OP_READ)
                    log("Connection accepted: $localAddress")
                }
            } else if (key.isReadable) {
                val client = key.channel() as SocketChannel
                val buffer = ByteBuffer.allocate(256)
                client.read(buffer)
                val result = buffer.data()
                log("Message received: $result")

                if (result == "Close") {
                    client.close()
                    log("Connection closed")
                }
            }
            iterator.remove()
        }
    }
}

fun stashLine() {
    stashed = console.getCursorBuffer().copy();
    try {
        console.getOutput().write("\u001b[1G\u001b[K");
        console.flush();
    } catch (e: IOException) {
        // ignore
    }
}

fun unstashLine() {
    try {
        console.resetPromptLine(console.getPrompt(), stashed.toString(), stashed!!.cursor)
    } catch (e: IOException) {
        // ignore
    }
}

fun println(output: String){
    stashLine()
    System.out.println(output)
    unstashLine()
}