package net.estinet.EstiConsole


import java.util.ArrayList
import java.util.Collections

/**
 * Usage:
 *  * String msg = Ansi.Red.and(Ansi.BgYellow).format("Hello %s", name)
 *  * String msg = Ansi.Blink.colorize("BOOM!")

 * Or, if you are adverse to that, you can use the constants directly:
 *  * String msg = new Ansi(Ansi.ITALIC, Ansi.GREEN).format("Green money")
 * Or, even:
 *  * String msg = Ansi.BLUE + "scientific"

 * NOTE: Nothing stops you from combining multiple FG colors or BG colors,
 * but only the last one will display.

 * @author dain
 */
class Ansi(vararg codes: String) {

    private val codes: Array<out String>
    private val codes_str: String

    init {
        this.codes = codes
        var _codes_str = ""
        for (code in codes) {
            _codes_str += code
        }
        codes_str = _codes_str
    }

    fun and(other: Ansi): Ansi {
        val both = ArrayList<String>()
        Collections.addAll(both, *codes)
        Collections.addAll(both, *other.codes)
        return Ansi(*both.toTypedArray())
    }

    fun colorize(original: String): String {
        return codes_str + original + SANE
    }

    fun format(template: String, vararg args: Any): String {
        return colorize(String.format(template, *args))
    }

    companion object {

        // Color code strings from:
        // http://www.topmudsites.com/forums/mud-coding/413-java-ansi.html
        val SANE = "\u001B[0m"

        val HIGH_INTENSITY = "\u001B[1m"
        val LOW_INTENSITY = "\u001B[2m"

        val ITALIC = "\u001B[3m"
        val UNDERLINE = "\u001B[4m"
        val BLINK = "\u001B[5m"
        val RAPID_BLINK = "\u001B[6m"
        val REVERSE_VIDEO = "\u001B[7m"
        val INVISIBLE_TEXT = "\u001B[8m"

        val BLACK = "\u001B[30m"
        val RED = "\u001B[31m"
        val GREEN = "\u001B[32m"
        val YELLOW = "\u001B[33m"
        val BLUE = "\u001B[34m"
        val MAGENTA = "\u001B[35m"
        val CYAN = "\u001B[36m"
        val WHITE = "\u001B[37m"

        val BACKGROUND_BLACK = "\u001B[40m"
        val BACKGROUND_RED = "\u001B[41m"
        val BACKGROUND_GREEN = "\u001B[42m"
        val BACKGROUND_YELLOW = "\u001B[43m"
        val BACKGROUND_BLUE = "\u001B[44m"
        val BACKGROUND_MAGENTA = "\u001B[45m"
        val BACKGROUND_CYAN = "\u001B[46m"
        val BACKGROUND_WHITE = "\u001B[47m"

        val HighIntensity = Ansi(HIGH_INTENSITY)
        val Bold = HighIntensity
        val LowIntensity = Ansi(LOW_INTENSITY)
        val Normal = LowIntensity

        val Italic = Ansi(ITALIC)
        val Underline = Ansi(UNDERLINE)
        val Blink = Ansi(BLINK)
        val RapidBlink = Ansi(RAPID_BLINK)

        val Black = Ansi(BLACK)
        val Red = Ansi(RED)
        val Green = Ansi(GREEN)
        val Yellow = Ansi(YELLOW)
        val Blue = Ansi(BLUE)
        val Magenta = Ansi(MAGENTA)
        val Cyan = Ansi(CYAN)
        val White = Ansi(WHITE)

        val BgBlack = Ansi(BACKGROUND_BLACK)
        val BgRed = Ansi(BACKGROUND_RED)
        val BgGreen = Ansi(BACKGROUND_GREEN)
        val BgYellow = Ansi(BACKGROUND_YELLOW)
        val BgBlue = Ansi(BACKGROUND_BLUE)
        val BgMagenta = Ansi(BACKGROUND_MAGENTA)
        val BgCyan = Ansi(BACKGROUND_CYAN)
        val BgWhite = Ansi(BACKGROUND_WHITE)
    }
}