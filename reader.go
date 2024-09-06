package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

func readProcUsage(PID string) (string, error) {
	fileName := "/proc/" + PID + "/stat"
	file, err := os.Open(fileName)

	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(file)

	line, _, err := reader.ReadLine()

	if err != nil {
		panic(err)
	}

	file.Close()
	return string(line), nil
}
func readProcCmdline(PID string, moreInfo bool) (string, error) {
	fileName := "/proc/" + PID + "/cmdline"
	file, err := os.Open(fileName)

	if err != nil {
		return "", err
	}
	reader := bufio.NewReader(file)

	line, _, err := reader.ReadLine()

	if err != nil {
		return "", err
	}

	file.Close()
	temp := string(line)
	if moreInfo {
		return temp, nil
	} else {
		splitted := strings.Split(temp, "-")
		return splitted[0], nil

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

			if x > 0.01 && err == nil {
				dataChan <- procUsage{y, file.Name(), x}
				//fmt.Printf(" Proccess: %s Usage: %.2f%% PID: %s \n", y, x, file.Name())
			}

		}
		time.Sleep(1 * time.Second)
	}

}
