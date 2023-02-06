package main

import (
	"pl.peth/serial-web-debugger/serial"
	"pl.peth/serial-web-debugger/web"
)

func main() {
	dataChannel := make(chan []byte)
	serial.DataChannel = dataChannel
	web.DataChannel = dataChannel

	serial.Init()
	web.Start()
}
