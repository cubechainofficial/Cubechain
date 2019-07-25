package main

import (
	"fmt"
	"time"
	"strings"
	"strconv"
	"encoding/gob"
	 "./config"
	 "./core"
)

var echo=fmt.Println
var Configure config.Configuration
var mstr="miningtesting!!..."
var	addr="CLQUKEdCeWmPzAmyJdHzo9cTBrq2JCBbPC"


func init() {
	Configure=config.LoadConfiguration("./config/cubechain.conf")
	core.Configure=Configure
	core.CubenoSet()
	echo (core.CubeSetNum)
	
	if core.GenFile=="" {
		path:="./config/genfile"
		core.GenFile=core.FileReadString(path)
		line:=strings.Split(core.GenFile,"\r\n")
		for _,v:=range line {
			result:=strings.Split(v, "|")
			genno,ok:=strconv.Atoi(result[0])
			if ok==nil {
				core.GenBlock[genno-1]+=v+"\r\n"
			}
		}
	}

	gob.Register(&core.TxData{})
	gob.Register(&core.TxBST{})
	gob.Register(map[string]string{})
}


func main() {
	go quickmining2()
	core.ServerRun()
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

func quickmining2() {
	tickChan:=time.Tick(2*time.Second)
	tickChan2:=time.Tick(time.Duration(Configure.Blocktime+1)*time.Second)
	echo("Cubechain start!")
	startTime:=time.Now()
    exeTime:=startTime.Add(-time.Duration(Configure.Blocktime)*time.Second)
    exeTime=exeTime.Add(-5*time.Second)
    cubeDuration:=time.Duration(Configure.Blocktime) * time.Second
	for {
		select {
		case <-tickChan:
			if  time.Since(exeTime)>=cubeDuration {
				exeTime=time.Now()
				cubemining2()
			}
		case <-tickChan2:
			go core.AllIndexing(0)
			go core.AllStatistic(0)
		}
	}
	echo("Cubechain end!")
}


func cubemining() { 
	var c core.Cube
	ch:=core.CubeHeight()+1
	echo (ch)
	c.Input(ch)
}

func cubemining2() { 
	var c core.Cube
	ch:=core.CubeHeight()
	ch2:=core.GetCubeHeight3()
	ch3,_:=strconv.Atoi(ch2)

	if ch3>ch {
		for i:=ch;i<=ch3;i++ {
			c.Cubeno=i
			c.CHash=""
			c.Download()
		}
	} else {
		if ch>3 {
			c.Cubeno=ch-2
			c.CHash=""
			c.Download()
			
			c.Cubeno=ch-1
			c.CHash=""
			c.Download()
		}

		echo (ch)
		c.InputChanel(ch)
	}
}
