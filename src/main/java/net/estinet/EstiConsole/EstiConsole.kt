package net.estinet.EstiConsole

import com.sun.xml.internal.fastinfoset.util.StringArray
import java.util.*


var version: String = "v0.0.1-BETA"
var mode: Modes = Modes.SPIGOT
var commands = ArrayList<ConsoleCommand>()

fun main(args: Array<String>) {
    println(Locale.getLocale(LocaleType.ENABLING))
    enable(args)
}

fun enable(args: Array<String>){
    /*
     * Startup Processes:
     */
    var isMode = false
    for(value in Modes.values()){
        if(args[0] == value.toString()){
            isMode = true
            mode = value
        }
    }
    if(isMode){
        println("Mode selected: $mode")
        println("Setting up Locale...")
        Locale.setupLocale()
        println("Setting up Configuration...")

        startCommandProcess()
    }
    else{
        println(Locale.getLocale(LocaleType.ERR_ON_START))
        println("[Error] Incorrect mode specified!")
        println("Exiting program...")
        System.exit(0)
    }
}

fun disable(){
    println(Locale.getLocale(LocaleType.DISABLING))
    /*
    * Disable Processes:
    */
    println(Locale.getLocale(LocaleType.DISABLED))
    System.exit(0)
}

fun startCommandProcess(){
    while(true){
        val input = System.console().readLine()
        val inputParsed = input.split(" ")
        if(inputParsed[0].toLowerCase() == "esticonsole" || inputParsed[0].toLowerCase() == "ec"){
            var foundValue = false
            for(cc in commands){
                if(cc.cName.toLowerCase() == inputParsed[1]){
                    val args = ArrayList<String>()
                    val i = 0
                    while(i < inputParsed.size){
                        if(i != 0 && i != 1) args.add(inputParsed[i])
                    }
                    cc.run(args)
                    foundValue = true
                    break
                }
            }
            if(!foundValue) println("") //TODO Must call command in java console
        }
    }
}

fun println(output: String){
    println("${Locale.getLocale(LocaleType.PREFIX)} $output")
}