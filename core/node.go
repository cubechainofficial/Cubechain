package core

import (
	"fmt"
	"os"
	"net"
	"net/http"
	"io"
	"io/ioutil"
	"strings"
	"strconv"
	"time"
)

const protocol = "tcp"
var nodes []NodeAddr

func IpCheck() []string {
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil
	}
	addrs=append(addrs,host)
	return addrs
}

func netError(err error) {
	if err != nil && err != io.EOF {
		fmt.Println("Network Error : ", err)
	}
}

func NodeRegister(mode string) string {
	arr:=IpCheck()
	reader :=strings.NewReader("cmode="+mode+"&mac="+arr[0]+"&ip="+arr[1]+"&hostname="+arr[2]+"&network="+Configure.Network+"&nettype="+Configure.Nettype+"&chaintype="+Configure.Chaintype+"&portnum="+strconv.Itoa(Configure.Port)+"&blocktime="+strconv.Itoa(Configure.Blocktime))
	request, _ := http.NewRequest("POST", "http://"+Configure.Mainserver+"/node/node_work.html", reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, _ := client.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}


func SearchNode() bool{
	rs:=false
	ii:=10
	for i:=1;i<=ii;i++ {
		nodeGet:=NodeRegister("get")	
		if nodeGet=="Nothing Node." {
			fmt.Printf("[ %s] Not found node.\n", time.Now())
			time.Sleep(30*time.Second)
		} else {
			line := strings.Split(nodeGet, "||")
			for k:=range line {
				result := strings.Split(line[k], "|")
				nodes=append(nodes,NodeAddr{result[0],result[1],result[2]})
				rs=true
			}
			i=ii
		}
	}
	return rs
}

func NodeListening() bool {
	myAddr:="localhost:"+strconv.Itoa(Configure.Port)
	addr, err := net.ResolveTCPAddr("tcp4", myAddr)
	netError(err)
	listener, err := net.ListenTCP(protocol, addr)
	netError(err)
	go func(l *net.TCPListener) {
		for {
			conn, err := l.AcceptTCP()
			netError(err)
			defer conn.Close()
		}
	}(listener)
	return true
}


func StartNode() {
	sbool:=SearchNode()
	if sbool {
		myAddr:="localhost:"+strconv.Itoa(Configure.Port)
		nodeNet,_:=net.Listen(protocol,myAddr)
		defer nodeNet.Close()
		for i,_:=range nodes {
			conn, _ := net.Dial(protocol, nodes[i].Ip+":"+nodes[i].Port)
			defer conn.Close()
		}
	} else {
		fmt.Println("Connection failed. Please try again later.")
	}

}

