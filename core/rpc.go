package core

import (
	"os"
    "encoding/json"
	"log"
	"fmt"
	"net/http"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	Callno		int			`json:"callno"`
	Com			string		`json:"com"`
	Vars		map[string]string	`json:"vars"`
	Rmsg		string		`json:"rmsg"`
}

type Response struct {
	Callno		int			`json:"callno"`
	Com			string		`json:"com"`
	Rmsg		string		`json:"rmsg"`
	Result		map[string]string		`json:"result"`
}

type ResResult struct {
	Result		map[string]string	`json:"result"`
}

type NodeConf struct {
	Network		string		`json:"network"`
	NodeID		string		`json:"node"`
	Status		string		`json:"status"`	
}

type PeerConf struct {
	PeerCount	string		`json:"peer_count"`
	MainIp		string		`json:"main_network_ip"`
	Sync		string		`json:"sync"`	
	CHeight		string		`json:"current_height"`	
}

type PowConf struct {
	Pows		string		`json:"Pow"`
	Hashrate	string		`json:"hashrate"`
	Address		string		`json:"address"`	
}

type BaseConf struct {
	Rpcver		string		`json:"rpcver"`
	Posstate	string		`json:"posstate"`
	Cheight		string		`json:"cheight"`	
}

var getError=Response{0,"Error","",map[string]string{"Message":"Please use post method"}}
var resonly ResResult
var resBool bool


type RpcFunc struct{}

func GetRequestVars(m map[string]string,key string) string {
	for k,v:= range m {
		if k==key {
			return v
		}
	}
	return ""
}

func (rf *RpcFunc) SyncApi(request Request,response *int) error {
	switch request.Com {
	case "sync_height":
		*response=CurrentHeight()
	default: 
		*response=0
	}
	return nil
}

func (rf *RpcFunc) Api(request Request, res *Response) error {
	var Peer PeerConf
	var Pow PowConf
	result:= make(map[string]string)
	resBool=false
	if GetRequestVars(request.Vars,"result_only")=="true" {
		resBool=true
	}
	switch request.Com {
	case "rpc_ver":
		result["ver"]="1"
	case "cube_height":
		result["cubeheight"]=strconv.Itoa(CurrentHeight())
	case "network_info":
		result["network"]=Configure.Network
		result["node"]=Configure.MainServer
		result["status"]=Configure.Nettype
	case "p2p_info":
		result["peer_count"]=Peer.PeerCount
		result["main_network_ip"]=Configure.MainServer
		result["sync"]=Peer.Sync
		result["current_height"]=Peer.CHeight
	case "cube_pow":
		result["pow"]=Pow.Pows
		result["hashrate"]=Pow.Hashrate
		result["address"]=Pow.Address
	case "cube_pos":
		result["pos"]="30000"
	case "cube_balance":
		addr:=GetRequestVars(request.Vars,"address")
		result["address"]=addr
		result["balance"]=strconv.FormatFloat(GetBalance(addr),'f',-1,64)
	case "cube_transaction_count":
		addr:=GetRequestVars(request.Vars,"address")
		result["address"]=addr
		result["txcount"]=strconv.Itoa(GetTransactionCount(addr))
	case "cube_transaction_list":
		addr:=GetRequestVars(request.Vars,"address")
		result["address"]=addr
		result["txlist"]=GetTransactionList(addr)
	case "cube_tx_list":
		addr:=GetRequestVars(request.Vars,"address")
		result["address"]=addr
		result["txlist"]=GetTxListDetail(addr)
	case "cube_info":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		cubeno,_:=strconv.Atoi(cubenos)
		result["info"]=GetCube(cubeno)
	case "block_info":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		blocknos:=GetRequestVars(request.Vars,"blockno")
		cubeno,_:=strconv.Atoi(cubenos)
		blockno,_:=strconv.Atoi(blocknos)
		result["info"]=GetBlock(cubeno,blockno)
	case "block_all":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		cubeno,_:=strconv.Atoi(cubenos)
		result["info"]=GetBlockAll(cubeno)
	case "cube_tx":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		cubeno,_:=strconv.Atoi(cubenos)
		result["tx"]=GetCubeTx(cubeno)
	case "block_tx":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		blocknos:=GetRequestVars(request.Vars,"blockno")
		cubeno,_:=strconv.Atoi(cubenos)
		blockno,_:=strconv.Atoi(blocknos)
		result["tx"]=GetBlockTx(cubeno,blockno)
	case "block_tx_search":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		blocknos:=GetRequestVars(request.Vars,"blockno")
		addr:=GetRequestVars(request.Vars,"addr")
		u_addr:=GetRequestVars(request.Vars,"u_addr")
		coins:=GetRequestVars(request.Vars,"coin")
		cubeno,_:=strconv.Atoi(cubenos)
		blockno,_:=strconv.Atoi(blocknos)
		coin,_:=strconv.ParseFloat(coins,64)
		result["tx"]=GetBlockTxHash(cubeno,blockno,addr,u_addr,coin)
	case "block_tx_hash":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		blocknos:=GetRequestVars(request.Vars,"blockno")
		addr:=GetRequestVars(request.Vars,"addr")
		u_addr:=GetRequestVars(request.Vars,"u_addr")
		coins:=GetRequestVars(request.Vars,"coin")
		cubeno,_:=strconv.Atoi(cubenos)
		blockno,_:=strconv.Atoi(blocknos)
		coin,_:=strconv.ParseFloat(coins,64)
		result["hash"]=GetBlockTxHash(cubeno,blockno,addr,u_addr,coin)
	case "indexing_data":
		result["indexing"]=GetIndexing()
	case "indexing_addr":
		addr:=GetRequestVars(request.Vars,"address")
		result["address"]=addr
		result["indexing"]=GetIndexingAddr(addr)
	case "statistic_data":
		result["statistic"]=GetStatistic()
	case "statistic_addr":
		addr:=GetRequestVars(request.Vars,"address")
		result["address"]=addr
		result["statistic"]=GetStatisticAddr(addr)
	case "statistic_rank":
		count:=GetRequestVars(request.Vars,"count")
		item:=GetRequestVars(request.Vars,"item")
		item2:=GetRequestVars(request.Vars,"item2")
		result["count"]=count
		result["item"]=item
		result["item2"]=item2
		icount,_:=strconv.Atoi(count)
		result["rank"]=GetStatisticRank(item,item2,icount)
	case "statistic_issue":
		result["issue"]=GetIssue()
	case "cube_transaction_detail":
		txhash:=GetRequestVars(request.Vars,"txhash")
		result["txhash"]=txhash
		tx,ci:=GetTransactionDetail(txhash)
		result["index"]=strconv.Itoa(ci.Index)
		result["cubeno"]=strconv.Itoa(ci.CubeNum)
		result["timestamp"]=strconv.Itoa(tx.Timestamp)
		result["from"]=tx.From
		result["to"]=tx.To
		result["amount"]=strconv.FormatFloat(tx.Amount,'f',-1,64)
		result["fee"]=strconv.FormatFloat(tx.Fee,'f',-1,64)
		result["hash"]=tx.Hash
		result["sign"]=tx.Sign
		result["nonce"]=strconv.Itoa(tx.Nonce)
	case "cube_transaction_data":
		txhash:=GetRequestVars(request.Vars,"txhash")
		tx,ci:=GetTransactionData(txhash)
		result["index"]=strconv.Itoa(ci.Index)
		result["cubeno"]=strconv.Itoa(ci.CubeNum)
		result["timestamp"]=strconv.Itoa(tx.Timestamp)
		result["from"]=tx.From
		result["to"]=tx.To
		result["amount"]=strconv.FormatFloat(tx.Amount,'f',-1,64)
		result["fee"]=strconv.FormatFloat(tx.Fee,'f',-1,64)
		result["hash"]=tx.Hash
		result["sign"]=tx.Sign
		result["nonce"]=strconv.Itoa(tx.Nonce)
	case "download_cube":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		cubeno,_:=strconv.Atoi(cubenos)
		result["file"]=FilePath(cubeno)+filepathSeparator+CubeFileName(cubeno)
		result["file"]=strings.Replace(result["file"], Configure.Datafolder, "", 1)
	case "download_block":
		cubenos:=GetRequestVars(request.Vars,"cubeno")
		blocknos:=GetRequestVars(request.Vars,"blockno")
		cubeno,_:=strconv.Atoi(cubenos)
		blockno,_:=strconv.Atoi(blocknos)
		result["file"]=FilePath(cubeno)+filepathSeparator+BlockName(cubeno,blockno)
		result["file"]=strings.Replace(result["file"], Configure.Datafolder, "", 1)
	default: 
		result["error"]="Command not found"
	}	
	res.Callno=request.Callno
	res.Com=request.Com
	res.Rmsg=request.Rmsg
	res.Result=result
	resonly.Result=res.Result
	if resBool {
		decho(resonly)
	} else {
		decho(res)
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	b,err:=json.MarshalIndent(v, "", "  ")
	if err!=nil {
		return err
	}
	w.Header().Set("Content-Type","application/json; charset=UTF-8")
	w.Write(b)
	return nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	var rf RpcFunc
	switch r.Method {
	case "GET":
		WriteJSON(w,getError)
		fmt.Println("Invalid get access")
	case "POST":
		d := json.NewDecoder(r.Body)
		p := &Request{}
		err := d.Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		result:=&Response{Callno:p.Callno,Com:p.Com,Rmsg:p.Rmsg}
		err=rf.Api(*p,result)

		if result.Result["file"]!="" {
			echo(result.Result["file"])
			http.Redirect(w,r,"/files/"+result.Result["file"],http.StatusSeeOther)
		} else {
			if resBool {
				WriteJSON(w,resonly)
			} else {
				WriteJSON(w,result)
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println("Invalid access")
	}
}


func Sync() {
	syncperiod:=10
	tickChan:=time.Tick(time.Duration(syncperiod)*time.Second)
	SyncOne()
	for {
		select {
		case <-tickChan:
			SyncOne()
		}
	}
}

func SyncOne() {
	req:=&Request{Callno:1,Com:"cube_height",Rmsg:"sync"}
	ch:=CurrentHeight()
	res:=ClientRun(req)
	if res>ch {
		SyncDownload(ch)
	}
}

func SyncDownload(idx int) {
	path:=FilePath(idx)
	url:="116.124.128.194/download/"+strconv.Itoa(idx)
	filep,_:=DownloadFile(path,url)
	if filep>"" {
 		echo("Download Cubefile.["+strconv.Itoa(idx)+"]")
		// Cube Convert file : Block, Cubing
	} else {
 		echo("Download failure.["+strconv.Itoa(idx)+"]")
	}
}


func cubeFileDown(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	path:=""
	cubenos:=r.Form.Get("cubeno")
	blocknos:=r.Form.Get("blockno")
	cubeno,_:=strconv.Atoi(cubenos)
	blockno,_:=strconv.Atoi(blocknos)
	if cubeno>0 && blockno>0 {
		path=FilePath(cubeno)+filepathSeparator+BlockName(cubeno,blockno)
		path=strings.Replace(path, Configure.Datafolder, "", 1)
		http.Redirect(w,r,"/files/"+path,http.StatusSeeOther)		
	} else if cubeno>0 {
		path=CubePath(cubeno)
		path=strings.Replace(path, Configure.Datafolder, "", 1)
		http.Redirect(w,r,"/files/"+path,http.StatusSeeOther)		
	} else {
 		echo("Invalid access.")
	}
}

func ClientRun(req *Request) int{
	client, err := net.Dial("tcp",Configure.Rpcip+":"+strconv.Itoa(Configure.Rpcport))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	c := jsonrpc.NewClient(client)
	var res int
	err = c.Call("RpcFunc.SyncApi", req, &res)
	if err != nil {
		log.Fatal("Rpc call:", err)
	}
	return res
}

func ClientRunFlag() {
	client, err := net.Dial("tcp",Configure.Rpcip+":"+strconv.Itoa(Configure.Rpcport))
	if err != nil {
		log.Fatal("dialing:", err)
	}
	c := jsonrpc.NewClient(client)
	req := &Request{}
	err=json.Unmarshal([]byte(os.Args[2]),&req)
	var res Response
	err = c.Call("RpcFunc.Api", req, &res)
	if err != nil {
		log.Fatal("Rpc call:", err)
	}
	var b []byte
	var jsonRes string
	if resBool {
		b,err=json.Marshal(resonly)
		jsonRes=ByteToStr(b)
	} else {
		b,err=json.Marshal(res)
		jsonRes=ByteToStr(b)
	}
	fmt.Println(jsonRes)
}

func RpcServer() {
	cal := new(RpcFunc)
	server := rpc.NewServer()
	server.Register(cal)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp",":"+strconv.Itoa(Configure.Rpcport))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Printf("Cubechain RPC Server Start.")
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			echo (conn)
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
	return
}

func HttpServer() {
	fs:=http.FileServer(http.Dir(Configure.Datafolder))
	http.Handle("/files/", http.StripPrefix("/files",fs))
	http.HandleFunc("/",apiHandler)
	http.HandleFunc("/download/",cubeFileDown)
	log.Printf("Cubechain Http RPC Start.:"+strconv.Itoa(Configure.Httpport))
	http.ListenAndServe(":"+strconv.Itoa(Configure.Httpport), nil)
}

func ServerRun() {
	go HttpServer()
	RpcServer()
}