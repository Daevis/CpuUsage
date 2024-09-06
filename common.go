package main

//#include <unistd.h>
// long hello(){
// 		sysconf(_SC_CLK_TCK);
// }
import "C"
import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func getCpuHz() float64 {
	var hertz float64 = float64(C.hello())
	return hertz
}
func convert(s string) int {
	value, _ := strconv.Atoi(s)
	return value
}

func readUptime() float64 {
	f, _ := os.Open("/proc/uptime")
	defer f.Close()

	reader := bufio.NewReader(f)
	line, _, _ := reader.ReadLine()
	temp := strings.Split(string(line), " ")
	strUptime := strings.TrimSuffix(temp[0], "\n")
	uptime, _ := strconv.ParseFloat(strUptime, 64)
	return uptime

}

type cpu struct {
	name    string
	user    int
	nice    int
	system  int
	idle    int
	iowait  int
	irg     int
	softirq int
}

type calcData struct {
	total       int
	prevTotal   int
	idle        int
	prevIdle    int
	nonIdle     int
	prevNonIdle int
}

type procUsage struct {
	name  string
	pid   string
	usage float64
}

type prevData struct {
	prevUptime    float64
	prevTotaltime float64
}
