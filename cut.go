package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	procNum, _ := exec.Command("nproc").Output()
	trimProc := strings.TrimSuffix(string(procNum), "\n")
	cores, _ := strconv.Atoi(trimProc)
	cores++ // first core is a sum of all cores

	dataChan := make(chan []string, 1)
	go reader(dataChan)
	result := make(chan float64, 1)
	go analyzer(dataChan, cores, result)

	procUsageChan := make(chan procUsage, 1)
	go getProcUsage(procUsageChan)
out:
	for {
		select {
		case <-sig:
			fmt.Println("\nCtrl-C called")
			break out
		case data := <-result:
			fmt.Printf("Average Cpu Usage: %.2f %% \n", data)

		case procUsage := <-procUsageChan:
			fmt.Printf(" Proccess: %s Usage: %.2f%% PID: %s \n",
				procUsage.name, procUsage.usage, procUsage.pid)

		}
	}
}
