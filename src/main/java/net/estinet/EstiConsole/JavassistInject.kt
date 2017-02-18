package net.estinet.EstiConsole

import java.security.*
import java.lang.instrument.*
import java.util.*
import javassist.*


class JavassistInject : ClassFileTransformer {

    @Throws(IllegalClassFormatException::class)
    override fun transform(loader: ClassLoader, className: String, redefiningClass: Class<*>, domain: ProtectionDomain, bytes: ByteArray): ByteArray {
        return transformClass(redefiningClass, bytes)
    }

    private fun transformClass(classToTransform: Class<*>, b: ByteArray): ByteArray {
        var b = b
        val pool = ClassPool.getDefault()
        var cl: CtClass? = null
        try {
            cl = pool.makeClass(java.io.ByteArrayInputStream(b))
            val methods = cl!!.declaredBehaviors
            for (i in methods.indices) {
                if (methods[i].isEmpty == false) {
                    changeMethod(methods[i])
                }
            }
            b = cl.toBytecode()
        } catch (e: Exception) {
            e.printStackTrace()
        } finally {
            if (cl != null) {
                cl.detach()
            }
        }
        return b
    }

    @Throws(NotFoundException::class, CannotCompileException::class)
    private fun changeMethod(method: CtBehavior) {
        if (method.name == "sendRawMessage") {
            method.setBody("System.out.println(message);")
        }
    }
}