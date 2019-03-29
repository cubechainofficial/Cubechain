package main

import (
	"fmt"
	"time"
	 "./config"
	 "./core"
)

var echo=fmt.Println
var Configure config.Configuration
var mstr="miningtesting!!..."

func init() {
	Configure=config.LoadConfiguration("./config/cubechain.conf")
	core.Configure=Configure
	core.CubenoSet()
	echo (core.CubeSetNum)
}


func main() {
	quickmining()
}


func quickmining() {
	tickChan:=time.Tick(time.Duration(Configure.Blocktime)*time.Second)
	echo("Cubechain start!")
	cubemining2()
	for {
		select {
		case <-tickChan:
			cubemining2()
		}
	}
	echo("Cubechain end!")
}


