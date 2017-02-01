package net.estinet.EstiConsole


var version: String = "0.0.1-BETA"
var mode: Modes = Modes.SPIGOT

fun main(args: Array<String>) {
    println("EstiConsole $version starting up...")
    enable(args);
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

    }
    else{
        println("[Error] Incorrect mode specified!")
        println("Exiting program...")
        System.exit(0)
    }
}

fun disable(){
    println("EstiConsole $version disabling...")
    /*
    * Disable Processes:
    */

}

class EstiConsole{

}