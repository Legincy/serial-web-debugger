package serial

import (
	"fmt"
	"go.bug.st/serial"
	"log"
	"strings"
)

var (
	DataChannel chan []byte
)

func Init() {
	portList, selection := requestTargetPort()

	if len(portList) == 0 {
		log.Fatal("No serial ports found!")
		return
	}

	selectedPort := portList[selection]

	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(selectedPort, mode)
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 1024)
	msgBuff := make([]byte, 0)
	go func() {
		for {
			bytesRead, err := port.Read(buff)
			if err != nil {
				log.Fatal(err)
				break
			}
			if bytesRead == 0 {
				fmt.Println("\nEOF")
				break
			}

			if DataChannel == nil {
				continue
			}

			msgBuff = append(msgBuff, buff[:bytesRead]...)
			if strings.Contains(string(msgBuff), "\n") {
				s := strings.SplitN(string(msgBuff), "\n", 2)
				msgBuff = []byte(s[1])
				DataChannel <- []byte(s[0])
			}
		}
	}()
}

func requestTargetPort() ([]string, int) {
	var selectedPort int

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Connected Ports ===")
	for index, port := range ports {
		fmt.Printf("[%d]\t%v\n", index+1, port)
	}

	if len(ports) > 0 {
		for {
			fmt.Printf("Select port which will be read (1-%d):\n", len(ports))
			fmt.Scan(&selectedPort)

			if selectedPort >= 1 && selectedPort <= len(ports) {
				selectedPort -= 1
				break
			}
		}
	}

	return ports, selectedPort
}
