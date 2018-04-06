package main

import (
	"os"
	"os/exec"
	"io"
	"bufio"
	"time"
	"io/ioutil"
	"strings"
	"strconv"
)

var Servers = make(map[string]*Server)

/*
 * Server struct and methods
 */

type Server struct {
	Settings   ServerConfig
	Log        []string
	Channel    chan string
	Process    *exec.Cmd
	OutputPipe io.ReadCloser
	ErrPipe    io.ReadCloser
	InputPipe  io.WriteCloser
	AutoStart  bool
	IsOnline   bool
}

//warning: this is a synchronous call.
func (server *Server) start() {
	info("Starting " + server.Settings.InstanceName)
	server.addLog("Starting " + server.Settings.InstanceName)
	strs := strings.Split(server.Settings.CommandToRun, " ")

	server.Process = exec.Command(strs[0], strs[1:]...) //doesn't execute the command just yet
	server.Process.Dir = server.Settings.HomeDirectory  //set working directory

	//Handle minecraft related tasks
	if server.Settings.MinecraftMode {
		if _, err := os.Stat(server.Settings.HomeDirectory + "/update"); os.IsNotExist(err) {
			os.Mkdir(server.Settings.HomeDirectory+"/update", 0755)
			info("Created the update directory for " + server.Settings.InstanceName + "!")
		}
		files, err := ioutil.ReadDir(server.Settings.HomeDirectory + "/update")
		if err != nil {
			server.addLog("[ERROR] Error reading directory " + server.Settings.HomeDirectory + "/update " + err.Error())
		}

		for _, f := range files { //move files from /update to plugins
			err := os.Rename(server.Settings.HomeDirectory+"/update/"+f.Name(), server.Settings.HomeDirectory+"/plugins/"+f.Name())
			if err != nil {
				server.addLog("[ERROR] Plugin update error: " + err.Error())
			}
			server.addLog("[INFO] Updated plugin " + f.Name())
		}
	}

	//Initializes input and output pipes
	server.OutputPipe, _ = server.Process.StdoutPipe()
	server.ErrPipe, _ = server.Process.StderrPipe()
	pipe, err := server.Process.StdinPipe()
	if err != nil {
		errMsg := "Process error for " + server.Settings.InstanceName + ": " + err.Error()
		info(errMsg)
		server.addLog(errMsg)
	}
	server.InputPipe = pipe

	//Start process
	err2 := server.Process.Start()
	if err2 != nil {
		errMsg := "Error starting process " + server.Settings.InstanceName + ": " + err2.Error()
		info(errMsg)
		server.addLog(errMsg)
		return
	}
	server.IsOnline = true

	//Function called when process ends
	deferFunc := func() {
		server.Process.Wait()
		server.IsOnline = false
		server.InputPipe.Close()
		server.OutputPipe.Close()
		info(server.Settings.InstanceName + " has stopped.")
		server.addLog(server.Settings.InstanceName + " has stopped.")
		if server.AutoStart {
			time.Sleep(time.Second * 2) //Let's not die right TODO
			go server.start()
		}
	}

	defer deferFunc()

	//Print output
	buff := bufio.NewScanner(server.OutputPipe)
	for buff.Scan() {
		server.addLog(buff.Text())
		if curServerView != nil && server.Settings.InstanceName == curServerView.Settings.InstanceName {
			println(buff.Text()) //prints reading from stdout
		}
	}
	buff2 := bufio.NewScanner(server.ErrPipe)
	for buff2.Scan() {
		server.addLog(buff2.Text())
		if curServerView != nil && server.Settings.InstanceName == curServerView.Settings.InstanceName {
			println(buff2.Text()) //prints reading from stdout
		}
	}
}

func (server *Server) kill() {
	server.AutoStart = false
	if err := server.Process.Process.Kill(); err != nil {
		info("Failed to kill process" + server.Settings.InstanceName + ": " + err.Error())
	} else {
		info("Killed " + server.Settings.InstanceName)
	}
}

func (server *Server) stop() {
	server.input(server.Settings.StopProcessCommand)
}

func (server *Server) input(input string) {
	io.WriteString(server.InputPipe, input+"\n")
}

func (server *Server) addLog(str string) {
	server.Log = append(server.Log, str)
	addToLogFile(str, logDirPath+"/"+server.Settings.InstanceName+"/current.log", logDirPath+"/"+server.Settings.InstanceName) //Write to log file
}

func (server *Server) getLog(beginIndex int, endIndex int) []string {
	if beginIndex < 0 {
		beginIndex = 0
	}
	return server.Log[beginIndex:endIndex]
}

func (server *Server) getLatestLogID() int {
	return len(server.Log) - 1
}

//END OF SERVER METHODS

//*************************************************************//

/*
 * Init and start all clients
 */

func ClientsStart() {
	for _, element := range instanceSettings.Servers {
		server := new(Server)
		server.Settings = element
		Servers[server.Settings.InstanceName] = server
		server.AutoStart = true

		info("Initialized server " + server.Settings.InstanceName + ".")

		go server.start()
	}
	info("Completed process initialization!")
}

/*
 * Stop all clients
 */

func ClientsStop() {
	info("Stopping all clients...")

	for key, _ := range Servers {
		if Servers[key].IsOnline {
			info("Stopping " + key + "...")
			go func(server *Server) {
				server.AutoStart = false
				server.stop()
				time.Sleep(time.Second * time.Duration(server.Settings.UnresponsiveKillTimeSeconds)) //perhaps remove thread blocking for the length of unresponsive kill time
				if server.IsOnline {
					server.Process.Process.Kill()
				}
			}(Servers[key])
		}
	}
}

/*
 * Kill clients
 */

func ClientsKill() {
	for key, _ := range Servers {
		if Servers[key].IsOnline {
			go func(server *Server) {
				server.AutoStart = false
				server.kill()
			}(Servers[key])
		}
	}
}

/*
 * Global helper functions for statically doing client manipulation
 */

func StartClient(name string) string {
	if _, ok := Servers[name]; ok {
		if Servers[name].IsOnline {
			return "Process already online."
		} else {
			Servers[name].AutoStart = true
			go Servers[name].start()
			return "Started " + Servers[name].Settings.InstanceName
		}
	} else {
		return "Server not found."
	}
}

/*
 * WARNING: Blocking thread
 */

func StopClient(name string) string {
	if _, ok := Servers[name]; ok {
		Servers[name].AutoStart = false
		if !Servers[name].IsOnline {
			return "Process already offline."
		} else {
			Servers[name].stop()
			for i := +0; uint(i) < Servers[name].Settings.UnresponsiveKillTimeSeconds; i++ {
				time.Sleep(time.Second)
				if !Servers[name].IsOnline {
					return "Stopped " + Servers[name].Settings.InstanceName
				}
			}
			Servers[name].kill()
			return "Didn't stop in time, killed " + Servers[name].Settings.InstanceName + " after waiting " + strconv.Itoa(int(Servers[name].Settings.UnresponsiveKillTimeSeconds)) + " seconds."
		}
	} else {
		return "Server not found."
	}
}

func KillClient(name string) string {
	if _, ok := Servers[name]; ok {
		Servers[name].AutoStart = false
		if !Servers[name].IsOnline {
			return "Process is not online."
		} else {
			Servers[name].kill()
			return "Killed process " + name + "."
		}
	} else {
		return "Server not found."
	}
}

func GetCPUUsage() string { /*
	if runtime.GOOS == "linux" {
		stat, err := linuxproc.ReadStat("/proc/stat")
		if err != nil {
			info("[ERROR] CPU stat read fail")
		}
		str := ""
		for i, s := range stat.CPUStats { //Loop through all cpu cores
			str += "CPU " + string(i) + ":\n"
			str += "User: " + string(s.User) + ", Nice: " + string(s.Nice) + ", System: " + string(s.System) + ", Idle: " + string(s.Idle) + ", IOWait: " + string(s.IOWait) + "\n"
		}
		return str
	}
	return "platform not supported" */
	return ""
}
func GetMemoryUsage() string {
	/*if runtime.GOOS == "linux" { TODO
		stat, err := linuxproc.ReadMemInfo("/proc/meminfo")
		if err != nil {
			info("[ERROR] Memory stat read fail")
		}
		str := ""
		str += "Free Memory: " + string(stat.MemFree) + "\nMemory Available: " + string(stat.MemAvailable) + "\nMemory Total: " + string(stat.MemTotal) + "\nSwap Free: " + string(stat.SwapFree) + "\nSwap Total: " + string(stat.SwapTotal) + "\n"
		return str
	}
	return "platform not supported"*/
	return ""
}
