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


func working() {
	var block core.Block
	block.Cubeno=651
	block.Blockno=2
	ph:=block.GetPrevHash()
	echo (ph)
}

func cubedown() {
	//var c core.Cube
	core.CubeDownload(600)

}

func cubemining() { 
	var c core.Cube
	ch:=core.CubeHeight()+1
	echo (ch)
	c.Input(ch)
}

func cubemining2() { 
	var c core.Cube
	ch:=core.CubeHeight()+1
	echo (ch)
	c.InputChanel(ch)
}


func blockmining1() {
	var b core.Block
	b.Input(1,10)
}
func blockmining2() { 
	var b core.Block
	c:=core.CubeHeight()+1
	echo (c)
	b.Input(c,10)
}




func blockscan() {
	b:=core.BlockScan(3042,3)
	b.Print()
}

func testcon() {
	r:=core.NodeSend2("pool_result","0&cubeno=1&blockno=3&hashstr=293u89u4832u48eaujfhugjnxcjnujdusifu")
	echo (r)
}

func pohcheck() {
	ph:=core.PohSet(1)
	echo (ph.Cubeno)
}

 
func checking() {
	r:=core.TxPool(7,3)
	bst,_:=core.TxpoolToBst(r)
	echo (r)
	bst.TreePrint2()

	b,_:=core.TxBlockData(7,3)
	echo (b)
}


