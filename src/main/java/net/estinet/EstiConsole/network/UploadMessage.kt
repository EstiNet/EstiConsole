package net.estinet.EstiConsole.network

import com.corundumstudio.socketio.AckRequest
import com.corundumstudio.socketio.SocketIOClient
import net.estinet.EstiConsole.EstiConsole
import java.io.*
import java.nio.file.Files
import java.nio.file.Paths
import java.util.*
import javax.xml.bind.DatatypeConverter
var uploadInProgress = ArrayList<String>()

class UploadMessage : Message{
    override val name: String = "upload"
    override fun run(args: List<String>, session: SocketIOClient, ack: AckRequest) {
        if(uploadInProgress.contains(session.remoteAddress.toString())){
            try{
                var f: File = File(args[0])
                var str = ""
                var i = 1
                while(i < args.size){
                    str += args[i] + " "
                    i++
                }
                var ind = DatatypeConverter.parseBase64Binary(str)
                val bw = BufferedOutputStream(FileOutputStream(f, true))
                bw.write(ind) //CHANGE THIS
                bw.flush()
                bw.close()

                ack.sendAckData("uploadcontinue")
            }
            catch(e: Throwable){
                e.printStackTrace()
                uploadInProgress.remove(session.remoteAddress.toString())
                ack.sendAckData("ecerror", "901")
            }
        }
        else{
            var f: File = File(args[0])
            try {
                if (f.exists()) {
                    if (EstiConsole.debug) {
                        EstiConsole.println("[Debug] upload request: Already exists. " + args[0])
                    }
                    ack.sendAckData("ecerror", "203")
                } else {
                    if (EstiConsole.debug) {
                        EstiConsole.println("[Debug] upload request: Writing to " + args[0])
                    }
                    f.createNewFile()
                    var str = ""
                    var i = 1
                    while(i < args.size){
                        str += args[i] + " "
                        i++
                    }
                    var ind = DatatypeConverter.parseBase64Binary(str)
                    Files.write(Paths.get(f.toURI()), ind) //CHANGE

                    if(EstiConsole.debug) {
                        EstiConsole.println("[Debug] upload request: Write complete.")
                    }
                    uploadInProgress.add(session.remoteAddress.toString())
                    ack.sendAckData("uploadcontinue", "uploadcontinue")
                }
            }
            catch(e: Throwable){
                e.printStackTrace()
                ack.sendAckData("ecerror", "901")
            }
        }
    }
}