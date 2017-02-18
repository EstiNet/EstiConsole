package net.estinet.EstiConsole

import javassist.*
import net.lingala.zip4j.core.ZipFile
import net.lingala.zip4j.exception.ZipException
import net.lingala.zip4j.model.ZipParameters
import net.lingala.zip4j.util.Zip4jConstants
import java.io.File
import java.io.FileOutputStream
import java.io.IOException
import java.util.jar.JarFile


object ASMInject{
    fun injectCode() {
        if (mode == Modes.PAPERSPIGOT) {
            //val className = "org.bukkit.craftbukkit.v1_11_R1.command.CraftConsoleCommandSender"
            val className = "org.bukkit.craftbukkit.v1_11_R1.command.CraftConsoleCommandSender"
            val jarFile = JarFile(File("cache").listFiles()[1])
            val f: File = File("cache").listFiles()[1]
// lets get a reference to the .class-file contained in the JAR
            val zipEntry = jarFile.getEntry(className.replace(".", "/") + ".class")
            if (zipEntry == null) {
                jarFile.close()
            }
// with our valid reference, we are now able to get the bytes out of the jar-archive
            val fis = jarFile.getInputStream(zipEntry!!)
            val classBytes = ByteArray(fis.available())
            fis.read(classBytes)
            val cp = ClassPool.getDefault()
            cp.insertClassPath(ClassClassPath(this.javaClass))
            var cp1: ClassPath? = null
            val cp2: ClassPath

// add the JAR file to the classpath
            try {
                cp1 = cp.insertClassPath(f.absolutePath)
            } catch (e1: NotFoundException) {
                e1.printStackTrace()
            }

// add the class file we are going to modify to the classpath
            cp2 = cp.appendClassPath(ByteArrayClassPath(className, classBytes))

            var modifiedBytes: ByteArray? = null
            try {
                val cc = cp.get(className)
                // skip instrumentation if the class is frozen and therefore
                // can't be modified
                if (!cc.isFrozen) {
                    cc.getDeclaredMethod("sendRawMessage").insertBefore("System.out.println(message);")
                }

                val destination = "cache/cache/"

                try {
                    val zipFile = ZipFile(f.absolutePath)
                    File(destination).mkdir()
                    zipFile.extractAll(destination)
                    File(destination + className.replace(".", "/") + ".class").delete()
                    File(destination + className.replace(".", "/") + ".class").createNewFile()
                    val s = FileOutputStream(File(destination + className.replace(".", "/") + ".class"))
                    s.write(cc.toBytecode())
                    s.close()
                    f.deleteRecursively()
                    zipFile.file.delete()
                    val parameters = ZipParameters()
                    parameters.compressionMethod = Zip4jConstants.COMP_DEFLATE
                    parameters.compressionLevel = Zip4jConstants.DEFLATE_LEVEL_NORMAL
                    zipFile.addFolder(destination, parameters)
                    File(destination).deleteRecursively()
                } catch (e: ZipException) {
                    e.printStackTrace()
                }

                //modifiedBytes = cc.toBytecode()
            } catch (e: NotFoundException) {
                e.printStackTrace()
            } catch (e: IOException) {
                e.printStackTrace()
            } catch (e: CannotCompileException) {
                e.printStackTrace()
            } catch (e: ClassNotFoundException) {
                e.printStackTrace()
            } finally {
                // free the locked resource files
                cp.removeClassPath(cp1!!)
                cp.removeClassPath(cp2!!)
                jarFile.close()
            }
        }
    }
}
