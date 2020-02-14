package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func main() {

	var cmdstr string
	var sampletime int
	if len(os.Args) >= 2 {
		fmt.Println(os.Args[1])
		cmdstr = os.Args[1]
	} else {
		fmt.Println("/usr/bin/yes")
		cmdstr = "/usr/bin/yes"

	}

	if len(os.Args) >= 3 {
		var err error
		sampletime, err = strconv.Atoi(os.Args[2])
		fmt.Println(sampletime)
		if err != nil {
			fmt.Println(err)
		}

	} else {
		sampletime = 2
	}

	//cmdstr := "/usr/bin/yes"
	done := make(chan bool)
	cmd := exec.Command(cmdstr)
	go func() {
		cmd.Run()
		mem := cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss
		cpuuser := cmd.ProcessState.SysUsage().(*syscall.Rusage).Utime
		cpusys := cmd.ProcessState.SysUsage().(*syscall.Rusage).Stime
		// cpupercent := float64((int64(cpu.Usec) / 1000)) / float64((int64(sampletime) * 1000)) * 100
		log.Printf("Memory: %vkB CPU: %v ms (user) %v ms (sys)", mem/1000, cpuuser.Usec/1000, cpusys.Usec/1000)
		done <- true
	}()
	for {
		select {
		case <-time.After(time.Duration(sampletime) * time.Second):
			log.Printf("Sampling completed, stopping process %v\n", cmdstr)
			err := cmd.Process.Signal(os.Interrupt)
			if err != nil {
				log.Fatalf("Can't stop process: %v\n", err)
			}
		case <-done:
			return
		}
	}
}
