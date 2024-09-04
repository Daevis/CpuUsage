package main

import (
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
			cpu := cpu{trimmed[0],
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
