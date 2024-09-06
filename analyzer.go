package main

import (
	"strconv"
	"strings"
	"time"
)

func analyzer(dataChan chan []string, cores int, result chan float64) {

	allCalcData := make([]calcData, cores)
	allCpuUsage := make([]float64, cores)
	cpus := make([]cpu, cores)
	for {
		data := <-dataChan
		avCpu := 0.0
		for i := range cores {
			trimmed := strings.Split(data[i], " ")
			cpu := cpu{
				trimmed[0],
				convert(trimmed[1]),
				convert(trimmed[2]),
				convert(trimmed[3]),
				convert(trimmed[4]),
				convert(trimmed[5]),
				convert(trimmed[6]),
				convert(trimmed[7])}
			cpus[i] = cpu
			idle := (cpu.idle + cpu.iowait)
			nonIdle := (cpu.user + cpu.nice + cpu.system + cpu.irg + cpu.softirq)
			total := idle + nonIdle
			temp := calcData{
				total,
				0,
				idle,
				0,
				nonIdle,
				0,
			}

			totald := temp.total - allCalcData[i].prevTotal
			idled := temp.idle - allCalcData[i].prevIdle
			var cpuUsage float64
			if totald != 0 {
				cpuUsage = (((float64(totald - idled)) * 100) / float64(totald))
			} else {
				cpuUsage = 0.0
			}
			allCpuUsage[i] = cpuUsage
			allCalcData[i] = temp
			allCalcData[i].prevIdle = allCalcData[i].idle
			allCalcData[i].prevTotal = allCalcData[i].total
			allCalcData[i].prevNonIdle = allCalcData[i].nonIdle

			if cpu.name != "cpu" {
				avCpu += cpuUsage
			}

		}
		avCpu /= float64(cores - 1)
		result <- avCpu
		time.Sleep(1 * time.Second)
	}
}

func getPIDUsage() func(PID string) (float64, string, error) {
	hashMap := make(map[string]*prevData)
	var hertz float64 = getCpuHz()

	return func(PID string) (float64, string, error) {
		if _, ok := hashMap[PID]; !ok {
			hashMap[PID] = &prevData{0.0, 0.0}
		}
		line, err := readProcUsage(PID)
		if err != nil {
			return 0, "", err
		}
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

		// Pass true for more info about proccess cmdline
		pidName, _ := readProcCmdline(PID, false)
		/*
			for i := 1; i < 10; i++ {
				if strings.Contains(splitted[i], ")") {
					pidName += splitted[i]
					break
				}
				pidName += splitted[i] + " "

			}
		*/
		hashMap[PID].prevTotaltime = totalTime
		hashMap[PID].prevUptime = uptime
		pidName = strings.Trim(pidName, "()")
		return cpuCurretn, pidName, nil
	}
}
