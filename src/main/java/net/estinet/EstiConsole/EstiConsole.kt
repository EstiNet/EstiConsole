package net.estinet.EstiConsole

import com.corundumstudio.socketio.SocketIOClient
import jline.console.ConsoleReader
import jline.console.CursorBuffer
import net.estinet.EstiConsole.commands.*
import net.estinet.EstiConsole.network.*
import org.fusesource.jansi.AnsiConsole
import java.io.File
import java.io.IOException
import java.io.OutputStreamWriter
import java.io.PrintWriter
import java.nio.file.Files
import java.util.*

object EstiConsole {
    var version: String = "v1.1.1"
    lateinit var javaProcess: Process
    lateinit var writer: PrintWriter
    var autoStartOnStop = false

    var debug = false;

    var logByteArray = ""

    /*
    * One of two println functions that must be used to print to console.
    */
    fun println(output: String) {
        stashLine()
        logByteArray += "\n${Locale.getLocale(LocaleType.PREFIX)} $output"
        SocketIO.sendToAll("log ${Locale.getLocale(LocaleType.PREFIX)} $output")
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
val commands = ArrayList<ConsoleCommand>()
val messages = ArrayList<Message>()

val sessions = HashMap<String, Boolean>();
val sessionStorage = HashMap<String, SocketIOClient>();

var networkOn = false;

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
 * Socket listener initializer
 */
fun setupMessages(){
    messages.add(HelloMessage())
    messages.add(CurlogsMessage())
    messages.add(CommandMessage())
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
        println("Starting network process...")
        Thread({ startNetworkServer()}).start()
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
    //EstiConsole.println("Injecting color code...")
    //ASMInject.injectCode()

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
        processCommand(input)
    }
}

fun processCommand(input: String){
    val inputParsed: List<String> = input.split(" ")
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

fun startNetworkServer(){
    SocketIO.doSocket()
}

fun parseJavaOutput(output: String) {

}

fun stashLine() {
    stashed = console.getCursorBuffer().copy();
    try {
        console.getOutput().write("\u001b[1G\u001b[K");
        console.flush();
    } catch (e: Exception) {
        // ignore
    }
}

fun unstashLine() {
    try {
        console.resetPromptLine(console.getPrompt(), stashed.toString(), stashed!!.cursor)
    } catch (e: Exception) {
        // ignore
    }
}
/*
 * One of two println functions that must be used to print to console.
 */
fun println(output: String){
    stashLine()
    EstiConsole.logByteArray += "\n$output"
    SocketIO.sendToAll("log $output")
    System.out.println(output)
    unstashLine()
}