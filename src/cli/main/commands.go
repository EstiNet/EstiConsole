package main

import (
	_ "log"
	"context"
	pb "../../protocol"
	"google.golang.org/grpc/connectivity"
	"time"
	"github.com/jroimartin/gocui"
	"fmt"
	"log"
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
	ch := make(chan int)
	go StartAttachSupervise(input, ch) //async supervisor
	<-ch
	attachCUI() //sync gui
}

var urgentCount uint64 = 30 //if the server should check more frequently for messages (message detection)

func StartAttachSupervise(input string, ch chan int) {
	ping := pb.ServerQuery{MessageId: -2, GetRam: true, GetCpu: true, ProcessName: input}

	ObtainNewLog(input, true) //initially fill slice
	ch <- 0
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
		attachLog = make([]string, reply2.MessageId+1) //fill initial with "" values
		attachLog = append(attachLog, reply2.Messages...)
		//writeslicetoview now in attachCli()
	} else {
		begin := len(attachLog) - 1 - int(reply2.MessageId)
		if begin < 0 {
			begin = 0
		}
		reply2.Messages = reply2.Messages[begin:len(reply2.Messages)]
		attachLog = append(attachLog, reply2.Messages...) //append new messages to log slice
		writeSliceToView(reply2.Messages, "v1")           //TODO only write screen height size
		prevBottomLine += len(reply2.Messages)
	}
}
func ObtainLogAtIndex(process string, index int) {
	obtain := pb.ServerQuery{MessageId: int64(index), GetRam: false, GetCpu: false, ProcessName: process}
	reply2, err2 := client.Attach(context.Background(), &obtain)
	checkError(err2) //caveat: can't accept 100 message gaps
	length := len(reply2.Messages)
	for i := 0; i < length; i++ {
		attachLog[index-i] = reply2.Messages[length-i-1]
	}
	(*cuiGUI).Update(func(g *gocui.Gui) error {
		out, err := (*cuiGUI).View("v1")
		if err != nil {
			log.Fatal(err)
		}
		out.Clear()
		for _, str := range attachLog {
			fmt.Fprintln(out, ""+str+"\u001b[0m")
		}
		return nil
	})
}
func UpdateInfo(cpu string, ram string) {
	(*cuiGUI).Update(func(g *gocui.Gui) error { //clear the view's text
		out, err := (*cuiGUI).View("v3")
		if err != nil {
			return err
		}
		out.Clear() //clear text
		fmt.Fprintln(out, cpu+"\n"+ram)
		return nil
	})
}
func SendCommand(command string, process string) {
	_, err := client.Attach(context.Background(), &pb.ServerQuery{MessageId: -2, Command: command, GetRam: false, GetCpu: false, ProcessName: process}) //initial ping
	checkError(err)
	urgentCount = 0
	ObtainNewLog(process, false)
}
