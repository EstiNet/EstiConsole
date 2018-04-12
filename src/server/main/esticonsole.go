package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"os/signal"
	"syscall"
	"time"
	"io"
	"log"
)

var (
	version          = "v2.0.0"
	instanceSettings InstanceConfig

	commands = make(map[string]interface{})

	curServerView *Server = nil

	logDirPath        = "./log"
	masterKey         string
	logCycle          = 0
	logDumpInProgress = false
	logQueue          = make([]LogAddition, 0)
)

type LogAddition struct {
	Str       string
	File      string
	Directory string
}

/*
 * Start the synchronous queue dump for logs
 * Note: please set logDumpInProgress to true before calling
 */

func startLogDump() {
	for !(len(logQueue) == 0) {
		addToLogFile(logQueue[0].Str, logQueue[0].File, logQueue[0].Directory)
		if len(logQueue) == 1 {
			logQueue = make([]LogAddition, 0)
		} else {
			logQueue = logQueue[1:]
		}
		time.Sleep(time.Millisecond * 1)
	}
	logDumpInProgress = false
}

/*
 * Output and logging related functions
 * TODO reduce chance of write collision with log buffer
 */

func addLog(str string) {
	if logCycle >= 20000 {
		InitLog()
		logCycle = 0
	}

	logQueue = append(logQueue, LogAddition{Str: str, File: logDirPath + "/current.log", Directory: logDirPath})

	if !logDumpInProgress {
		logDumpInProgress = true
		go startLogDump()
	}

	logCycle++
}
func addToLogFile(str string, file string, directory string) {
	//check if logdirpath exists
	if _, err := os.Stat(logDirPath); os.IsNotExist(err) {
		fmt.Println("The entire config directory is gone! WTH are you doing?\nCommencing repair.")
		InitLog()
		PostInitLog()
	}
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		fmt.Println("Logging directory " + directory + " disappeared! WAT\nCommencing repair.")
		os.Mkdir(directory, 0755)
		InitLogFile(directory)
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
	if os.IsNotExist(err) {
		fmt.Println("Config file disappeared! WTH are you doing?\nCommencing repair.")
		InitLog()
		PostInitLog()
		f, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0600)
	} else if err != nil {
		ClientsKill()
		log.Fatal(err)
	}

	defer f.Close()

	if _, err = f.WriteString(str + "\n"); err != nil {
		ClientsKill()
		log.Fatal(err)
	}
}
func logFatal(err error) {
	if err != nil {
		addLog(err.Error())
		ClientsKill()

		for i := 0; i < 2000 && logDumpInProgress; i++ { //wait for log dump to finish and exit
			time.Sleep(time.Millisecond)
		}
		log.Fatal(err)
	}
}
func logFatalStr(str string) {
	addLog(str)
	ClientsKill()
	for i := 0; i < 2000 && logDumpInProgress; i++ { //wait for log dump to finish and exit
		time.Sleep(time.Millisecond)
	}
	log.Fatal(str)
}
func println(str string) {
	addLog(str)
	fmt.Println(str)
}
func print(str string) {
	addLog(str)
	fmt.Print(str)
}
func info(str string) {
	addLog(str)
	println(time.Now().Format("2006-01-02 15:04:05") + " [INFO] " + str)
}
func debug(str string) {
	fmt.Println(str)
}

/*
 * Program operation related functions
 */

func init() {
	commands["version"] = CommandVersion
	commands["stop"] = CommandStop
	commands["help"] = CommandHelp
	commands["list"] = CommandList
	commands["switch"] = CommandSwitch
	commands["instancestop"] = CommandInstanceStop
	commands["start"] = CommandStart
	commands["kill"] = CommandKill
}

/*
 * Entry point for program.
 */

func main() {

	//Start logging
	InitLog()

	println("EstiConsole " + version + "")

	//System signal hooks
	info("Registering hooks...")
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		info("Received host " + sig.String())
		done <- true
	}()
	info("Completed!")

	//Continue with startup
	info("Setting up and loading configuration...")
	LoadConfig()
	info("Completed!")

	info("Completing post log initialization...")
	PostInitLog()

	info("Starting network processes...")
	go NetworkStart()
	info("Starting client processes...")
	go ClientsStart()
	info("Starting command system...")
	go ConsoleStart()

	//Receive interrupt
	<-done
	Shutdown()
}

/*
 * Shutdown task
 */
func Shutdown() {
	info("Commencing instance shutdown.")
	go ClientsStop()

	var maxKillTime uint
	for _, server := range Servers { //get the longest unresponsive kill time period
		if server.Settings.UnresponsiveKillTimeSeconds > maxKillTime {
			maxKillTime = server.Settings.UnresponsiveKillTimeSeconds
		}
	}

	for i := 0; uint(i) <= maxKillTime; i++ { //loop until the server is forced to shut down or all the processes have shut down
		if uint(i) == maxKillTime { //if servers are still online after server unresponsive kill time
			info("Force shutting down, processes still online.")
			break
		}

		stillOnline := false
		time.Sleep(time.Second)

		for _, server := range Servers {
			if server.IsOnline {
				stillOnline = true
				break
			}
		}
		if !stillOnline {
			break
		}
	}
	grpcServer.Stop()

	info("Exited EstiConsole " + version)
	os.Exit(0)
}

/*
 * Command handler
 */

func ConsoleStart() {
	var reader = bufio.NewReader(os.Stdin)
	for { //command line loop
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if err != nil {
			println(err.Error())
		}
		input = strings.TrimRight(input, "\n")
		cFound := false
		if strings.Split(input, " ")[0] == "ec" {
			for k, v := range commands {
				if k == strings.Split(input, " ")[1] {
					in := ""
					for i, str := range strings.Split(input, " ") {
						if i != 0 && i != 1 {
							in += str
						}
					}
					v.(func(string))(in)
					cFound = true
					break
				}
			}
			if !cFound {
				println("Unknown EstiConsole command.")
			}
		} else if curServerView == nil {
			println("No current server view. Please do /ec switch [server].")
		} else {
			curServerView.input(input)
		}
	}
}

func substring(s string, start int, end int) string {
	start_str_idx := 0
	i := 0
	for j := range s {
		if i == start {
			start_str_idx = j
		}
		if i == end {
			return s[start_str_idx:j]
		}
		i++
	}
	return s[start_str_idx:]
}
