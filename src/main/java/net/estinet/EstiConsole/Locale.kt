package net.estinet.EstiConsole

import java.io.*
import java.util.*

object Locale{
    var phrases = ArrayList<String>()
    val default = ArrayList<String>()
    fun load(){
        default.add("Unknown command! Do /ec help for help!")
        default.add("[EstiConsole]")
        default.add("EstiConsole \${EstiConsole.version} turning off...")
        default.add("EstiConsole \${EstiConsole.version} is off. Goodbye!")
        default.add("EstiConsole \${EstiConsole.version} starting up...")
        default.add("EstiConsole \${EstiConsole.version} is online!")
        default.add("Error while starting EstiConsole \${EstiConsole.version}")
    }
    fun setupLocale(){
        load()
        val f: File = File("languages.properties")
        var output: OutputStream? = null
        try{
            if(!f.exists()) f.createNewFile()
            val prop: Properties = Properties()
            val input: InputStream = FileInputStream(f)
            prop.load(input)
            output = FileOutputStream(f)
            for(type in LocaleType.values()){
                if (prop.getProperty(type.toString().toLowerCase()) == null)
                    prop.setProperty(type.toString().toLowerCase(), default.get(type.ordinal))
                phrases.add(prop.getProperty(type.toString().toLowerCase()))
            }
        prop.store(output, null)
        } catch (io: IOException) {
            io.printStackTrace();
        } finally {
            if (output != null) {
                try {
                    output.close();
                } catch (e: IOException) {
                    e.printStackTrace();
                }
            }
        }
    }
    fun getLocale(type: LocaleType): String{
        return phrases.get(type.ordinal).replace("\${EstiConsole.version}", EstiConsole.version)
    }
}
