package main

import (
	_ "log"
	"context"
	pb "../../protocol"
	"google.golang.org/grpc/connectivity"
	"time"
	"github.com/jroimartin/gocui"
	"fmt"
)

func CommandHelp(input string) {
	println("-----Help-----")
	println("-h               | Get the help interface for flags.")
	println("version          | Get the version of this instance.")
	println("status           | Get the status of the instance")
	println("instancestop     | Stop this instance of EstiConsole.")
	println("list             | List all of the client servers.")
	println("attach [process] | Switch view to another process.")
	println("stop [process]   | Stop the process using the default stop command.")
	println("start [process]  | Start the process.")
	println("kill [process]   | Forcibly kill the process.")
}

func CommandVersion(input string) {
	startCon()
	reply, err := client.Version(context.Background(), &pb.String{Str: "test reply"})
	checkError(err)
	println("Version: ", reply.Str)
}

func CommandList(input string) {
	startCon()
	reply, err := client.List(context.Background(), &pb.String{Str: ""})
	checkError(err)
	println(reply.Str)
}

func CommandStop(input string) {
	startCon()
	reply, err := client.Stop(context.Background(), &pb.String{Str: input})
	checkError(err)
	println(reply.Str)
}

func CommandInstanceStop(input string) {
	startCon()
	reply, err := client.InstanceStop(context.Background(), &pb.String{Str: ""})
	checkError(err)
	println(reply.Str)
}

func CommandStart(input string) {
	startCon()
	reply, err := client.Start(context.Background(), &pb.String{Str: input})
	checkError(err)
	println(reply.Str)
}

func CommandKill(input string) {
	startCon()
	reply, err := client.Kill(context.Background(), &pb.String{Str: input})
	checkError(err)
	println(reply.Str)
}

func CommandStatus(input string) {
	startCon()
	if conn.GetState() == connectivity.Ready {
		println("Connection successful!")
	}
}

/*
 * Attach + CUI related code
 */

var cpuInfo, ramInfo, procName string
var attachLog []string

func CommandAttach(input string) {
	procName = input

	startCon()
	go StartAttachSupervise(input) //async supervisor TODO FIX ASYNC load
	attachCUI() //sync gui
}

var urgentCount uint64 = 30 //if the server should check more frequently for messages (message detection)

func StartAttachSupervise(input string) {
	ping := pb.ServerQuery{MessageId: -2, GetRam: true, GetCpu: true, ProcessName: input}

	ObtainNewLog(input, true) //initially fill slice
	for {
		if urgentCount < 30 { //increase urgent count if there are periods of no messages (4 seconds)
			urgentCount++
		}
		reply, err := client.Attach(context.Background(), &ping) //initial ping
		checkError(err)

		if int(reply.MessageId) >= len(attachLog)-1 { //if there are new messages
			ObtainNewLog(input, false)
			urgentCount = 0
		}

		//UpdateInfo(reply.CpuUsage, reply.RamUsage) TODO cpu and ram usage

		if urgentCount >= 30 { //slow check for messages
			t, _ := time.ParseDuration("10ms")
			for i := 0; i < 130; i++ { //sleep 1300ms
				if urgentCount == 0 { //leave loop early if new messages are sent
					break
				}
				time.Sleep(t) //sleep 10ms
			}
		} else { //burst message detection
			t, _ := time.ParseDuration("100ms")
			time.Sleep(t)
		}
	}
}

func ObtainNewLog(process string, firstGet bool) {
	obtainNewest := pb.ServerQuery{MessageId: -1, GetRam: false, GetCpu: false, ProcessName: process}
	reply2, err2 := client.Attach(context.Background(), &obtainNewest)
	checkError(err2) //caveat: can't accept 100 message gaps
	if firstGet {
		if reply2.MessageId == 0 {
			reply2.MessageId++
		}
		attachLog = make([]string, reply2.MessageId-1)    //fill initial with "" values
		attachLog = append(attachLog, reply2.Messages...) //TODO duplication of previous message
		for _, cur := range attachLog {
			//println(cur)
			//writeToView("\033[30;1m" + cur + "\033[0m", "v1")
			writeToView(cur, "v1")
		}
	} else {
		reply2.Messages = reply2.Messages[(len(attachLog) - 1 - int(reply2.MessageId)):len(reply2.Messages)]
		attachLog = append(attachLog, reply2.Messages...) //append new messages to log slice
		for _, cur := range reply2.Messages {
			//println(cur)
			//writeToView("\033[30;1m" + cur + "\033[0m", "v1")
			writeToView(cur, "v1")
		}
	}
}
func ObtainLogAtIndex(process string, index int) {

}
func UpdateInfo(cpu string, ram string) {
	(**cuiGUI).Update(func(g *gocui.Gui) error { //clear the view's text
		out, err := (**cuiGUI).View("v3")
		if err != nil {
			return err
		}
		out.Clear()         //clear text
		fmt.Fprintln(out, cpu + "\n" + ram)
		return nil
	})
}
func SendCommand(command string, process string) {
	_, err := client.Attach(context.Background(), &pb.ServerQuery{MessageId: -2, Command: command, GetRam: false, GetCpu: false, ProcessName: process}) //initial ping
	checkError(err)
	urgentCount = 0
	ObtainNewLog(process, false)
}

