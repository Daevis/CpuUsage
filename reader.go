package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

//#include <unistd.h>
// long hello(){
// 		sysconf(_SC_CLK_TCK);
// }
import "C"

type procUsage struct {
	name  string
	pid   string
	usage float64
}

type prevData struct {
	prevUptime    float64
	prevTotaltime float64
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
func getPIDUsage() func(PID string) (float64, string, error) {
	hashMap := make(map[string]*prevData)
	var hertz float64 = float64(C.hello())

	return func(PID string) (float64, string, error) {
		if _, ok := hashMap[PID]; !ok {
			hashMap[PID] = &prevData{0.0, 0.0}
		}

		fileName := "/proc/" + PID + "/stat"
		file, err := os.Open(fileName)

		if err != nil {
			return 0.0, "", err
		}
		reader := bufio.NewReader(file)

		line, _, err := reader.ReadLine()

		if err != nil {
			panic(err)
		}

		file.Close()

		temp := string(line)
		splitted := strings.Split(temp, " ")
		utime, _ := strconv.Atoi(splitted[15])
		stime, _ := strconv.Atoi(splitted[16])
		cutime, _ := strconv.Atoi(splitted[17])
		cstime, _ := strconv.Atoi(splitted[18])
		totalTime := float64(utime + stime + cutime + cstime)
		uptime := readUptime()
		elapesedTime := (uptime - hashMap[PID].prevUptime)
		cpuCurretn := 100 * (((totalTime - hashMap[PID].prevTotaltime) / hertz) / elapesedTime)

		var pidName string
		for i := 1; i < 10; i++ {
			if strings.Contains(splitted[i], ")") {
				pidName += splitted[i]
				break
			}
			pidName += splitted[i] + " "

		}
		hashMap[PID].prevTotaltime = totalTime
		hashMap[PID].prevUptime = uptime
		pidName = strings.Trim(pidName, "()")
		return cpuCurretn, pidName, nil
	}
}
func reader(dataChan chan []string) {
	for {
		file, err := os.Open("/proc/stat")

		if err != nil {
			os.Exit(0)
		}
		reader := bufio.NewReader(file)
		var data [][]byte

		iter := 0
		for {
			line, _, err := reader.ReadLine()

			if len(line) > 0 {
				data = append(data, line)
			}

			if err != nil {
				break
			}
			iter++
		}
		file.Close()
		var temp []string
		for i := range data {
			temp = append(temp, string(data[i]))
		}
		dataChan <- temp
		time.Sleep(1 * time.Second)
		//fmt.Printf("data: %s \n lines: %d\n", data, iter)
	}
}
func getProcUsage(dataChan chan procUsage) {

	files, err := os.ReadDir("/proc/")
	if err != nil {
		panic(err)
	}
	f := getPIDUsage()

	for {

		for _, file := range files {
			_, err := strconv.Atoi(file.Name())
			if err != nil {
				continue
			}

			x, y, err := f(file.Name())

			if x > 0 && err == nil {
				dataChan <- procUsage{y, file.Name(), x}
				//fmt.Printf(" Proccess: %s Usage: %.2f%% PID: %s \n", y, x, file.Name())
			}

		}
		time.Sleep(1 * time.Second)
	}

}
