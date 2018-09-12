package core

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Version = "0.93"
var nodes []NodeAddr
var handelStr string
var handelData interface{}
var cubechainInfo Cubechain

type HandleFunc func(*bufio.ReadWriter)

type Nodepoint struct {
	listener net.Listener
	handler  map[string]HandleFunc
	m        sync.RWMutex
}

type complexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *complexData
}

func netError(err error) {
	if err != nil && err != io.EOF {
		log.Print("Network Error : ", err)
	}
}

func IpCheck() []string {
	host, err := os.Hostname()
	if err != nil {
		return nil
	}
	addrs, err := net.LookupHost(host)
	if err != nil {
		return nil
	}
	addrs = append(addrs, host)
	return addrs
}

func TcpDial(addr string) (*bufio.ReadWriter, error) {
	log.Print("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}

func NewNodepoint() *Nodepoint {
	nodeGet := NodeRegister("node", "get")
	log.Print("Node get : " + nodeGet)
	return &Nodepoint{
		handler: map[string]HandleFunc{},
	}
}

func (np *Nodepoint) AddHandleFunc(name string, f HandleFunc) {
	np.m.Lock()
	np.handler[name] = f
	np.m.Unlock()
}

func (np *Nodepoint) Listen() error {
	var err error
	np.listener, err = net.Listen("tcp", ":"+strconv.Itoa(Configure.Port))
	if err != nil {
		return err
	}
	log.Println("Listen on", np.listener.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := np.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go np.handleMessages(conn)
	}
}

func (np *Nodepoint) handleMessages(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()
	for {
		log.Print("Receive command '")
		cmd, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			log.Println("Reached EOF - close this connection.\n   ---")
			return
		case err != nil:
			log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
			return
		}
		cmd = strings.Trim(cmd, "\n ")
		log.Println(cmd + "'")
		np.m.RLock()
		handleCommand, ok := np.handler[cmd]
		np.m.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' is not registered.")
			return
		}
		handleCommand(rw)
	}
}

func NodeServer() error {
	Nodepoint := NewNodepoint()
	Nodepoint.AddHandleFunc("NODECHECK", handleNodeCheck)
	Nodepoint.AddHandleFunc("SYNC", handleSync)
	Nodepoint.AddHandleFunc("VERSION", handleVersion)
	Nodepoint.AddHandleFunc("DOWNLOAD", handleDownload)
	Nodepoint.AddHandleFunc("STRING", handleString)
	Nodepoint.AddHandleFunc("DATA", handleData)
	return Nodepoint.Listen()
}

func NodesGet(nodeGet string) []NodeAddr {
	line := strings.Split(nodeGet, "||")
	for k := range line {
		result := strings.Split(line[k], "|")
		nodes = append(nodes, NodeAddr{result[0], result[1], result[2]})
	}
	return nodes
}

func handleNodeCheck(rw *bufio.ReadWriter) {
	arr := IpCheck()
	reader := strings.NewReader("cmode=get&mac=" + arr[0] + "&ip=" + arr[1] + "&hostname=" + arr[2] + "&network=" + Configure.Network + "&nettype=" + Configure.Nettype + "&chaintype=" + Configure.Chaintype + "&portnum=" + strconv.Itoa(Configure.Port) + "&blocktime=" + strconv.Itoa(Configure.Blocktime))
	request, _ := http.NewRequest("POST", "http://"+Configure.Mainserver+"/node/node_work.html", reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, _ := client.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	s := string(body)
	r := NodesGet(s)
	log.Println(r)
}

func handleDownload(rw *bufio.ReadWriter) {
	log.Println("Download file")
	filename := chainFileName()
	fileUrl := "http://" + Configure.Mainserver + "/node/download.html?file=" + filename
	err := DownloadFile("./bdata/"+filename, fileUrl)
	if err != nil {
		log.Println(err)
	}
}

func handleString(rw *bufio.ReadWriter) {
	log.Print(handelStr)
	err := rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

func handleData(rw *bufio.ReadWriter) {
	log.Print("Receive data:")
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&handelData)
	if err != nil {
		log.Println("Error decoding data:", err)
		return
	}
	log.Printf("Outer data: \n%#v\n", handelData)
}

func handleSync(rw *bufio.ReadWriter) {
	handleVersion(rw)
	handleDownload(rw)
	cubechainInfo = ChainFileRead()
	handelStr = "Download chain file"
	handleString(rw)
	ccnt := ChainCheck(&cubechainInfo)
	handelStr = "Chain check : " + strconv.Itoa(ccnt)
	handleString(rw)
}

func handleVersion(rw *bufio.ReadWriter) {
	log.Println("Version Informaition:" + Version)
	err := rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func NodeRegister(comm string, mode string) string {
	arr := IpCheck()
	reader := strings.NewReader("cmode=" + mode + "&mac=" + arr[0] + "&ip=" + arr[1] + "&hostname=" + arr[2] + "&network=" + Configure.Network + "&nettype=" + Configure.Nettype + "&chaintype=" + Configure.Chaintype + "&portnum=" + strconv.Itoa(Configure.Port) + "&blocktime=" + strconv.Itoa(Configure.Blocktime))
	request, _ := http.NewRequest("POST", "http://"+Configure.Mainserver+"/chain/"+comm, reader)
	request.Header.Add("content-type", "application/x-www-form-urlencoded")
	request.Header.Add("cache-control", "no-cache")
	client := &http.Client{}
	res, _ := client.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body)
}

func SearchNode() bool {
	rs := false
	ii := 10
	for i := 1; i <= ii; i++ {
		nodeGet := NodeRegister("node", "get")
		if nodeGet == "Nothing Node." {
			fmt.Printf("[ %s] Not found node.\n", time.Now())
			time.Sleep(30 * time.Second)
		} else {
			line := strings.Split(nodeGet, "||")
			for k := range line {
				result := strings.Split(line[k], "|")
				nodes = append(nodes, NodeAddr{result[0], result[1], result[2]})
				rs = true
			}
			i = ii
		}
	}
	return rs
}

func NodeListening() string {
	cb := make(chan string)
	sb := "Node Listening start"
	arr := IpCheck()
	myAddr := arr[1] + ":" + strconv.Itoa(Configure.Port)
	addr, err := net.ResolveTCPAddr("tcp4", myAddr)
	listener, err := net.ListenTCP("tcp", addr)
	netError(err)
	go func(l *net.TCPListener) {
		for {
			connection, err := l.AcceptTCP()
			netError(err)
			sb = "connection" + strconv.Itoa(Configure.Port) + strconv.Itoa(int(time.Now().Unix()))
			cb <- sb
			fmt.Println(connection)
		}
	}(listener)
	return ""
}
