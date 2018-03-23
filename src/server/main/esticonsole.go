package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"os/signal"
	"syscall"
	"os/exec"
	"runtime"
	"time"
	"io"
	"log"
)

var version = "v2.0.1"
var instanceSettings InstanceConfig

var commands = make(map[string]interface{})

var curServerView *Server = nil

var logDirPath = "./log"

var clear map[string]func()

/*
 * Output and logging related functions
 * TODO make async with log queue
 */

func addLog(str string) {
	addToLogFile(str, logDirPath + "/current.log")
}
func addToLogFile(str string, directory string) {
	f, err := os.OpenFile(directory, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(str + "\n"); err != nil {
		panic(err)
	}
}
func logFatal(err error) {
	addLog(err.Error())
	log.Fatal(err)
}
func logFatalStr(str string) {
	addLog(str)
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

	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

/*
 * Entry point for program.
 */

func main() {

	//Start logging
	InitLog()

	println("EstiConsole " + version)

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
	ClientsStop()

	//TODO REPLACE WITH THREAD BLOCKING
	//TODO TEMP SOLUTION SHOULD NOT BE WORKING IN PRODUCTION!!!!!!
	time.Sleep(time.Second * 8)

	grpcServer.Stop()

	info("Exited EstiConsole " + version)
	os.Exit(0)
}

/*
 * Command handler
 */

func ConsoleStart() {
	var reader = bufio.NewReader(os.Stdin)
	for {
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

func ClearTerminal() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
