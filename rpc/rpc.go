package rpc

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
	//"strings"
	"../lib"
	"../config"
	"../core"
)

var Configure config.Configuration

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
		result["cubeheight"]=strconv.Itoa(core.CurrentHeight())
	case "network_info":
		result["network"]=Configure.Network
		result["node"]=Configure.Mainserver
		result["status"]=Configure.Nettype
	case "p2p_info":
		result["peer_count"]=Peer.PeerCount
		result["main_network_ip"]=Configure.Mainserver
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
		result["balance"]=strconv.Itoa(core.GetBalance(addr))
	case "cube_transaction_count":
		addr:=GetRequestVars(request.Vars,"address")
		result["address"]=addr
		result["txcount"]=strconv.Itoa(core.GetTransactionCount(addr))
	case "cube_transaction_list":
		addr:=GetRequestVars(request.Vars,"address")
		result["address"]=addr
		result["txlist"]=core.GetTransactionList(addr)
	case "cube_transaction_detail":
		txhash:=GetRequestVars(request.Vars,"txhash")
		result["txhash"]=txhash
		tx,ci:=core.GetTransactionDetail(txhash)
		result["index"]=strconv.Itoa(ci.Index)
		result["cubeno"]=strconv.Itoa(ci.CubeNum)
		result["timestamp"]=strconv.Itoa(tx.Timestamp)
		result["from"]=lib.ByteToStr(tx.From)
		result["to"]=lib.ByteToStr(tx.To)
		result["amount"]=strconv.Itoa(tx.Amount)
		result["fee"]=strconv.Itoa(tx.Fee)
		result["hash"]=lib.ByteToStr(tx.Hash)
		result["sign"]=lib.ByteToStr(tx.Sign)
		result["nonce"]=strconv.Itoa(tx.Nonce)
	case "cube_transaction_data":
		txhash:=GetRequestVars(request.Vars,"txhash")
		tx,ci:=core.GetTransactionData(txhash)
		result["index"]=strconv.Itoa(ci.Index)
		result["cubeno"]=strconv.Itoa(ci.CubeNum)
		result["timestamp"]=strconv.Itoa(tx.Timestamp)
		result["from"]=lib.ByteToStr(tx.From)
		result["to"]=lib.ByteToStr(tx.To)
		result["amount"]=strconv.Itoa(tx.Amount)
		result["fee"]=strconv.Itoa(tx.Fee)
		result["hash"]=lib.ByteToStr(tx.Hash)
		result["sign"]=lib.ByteToStr(tx.Sign)
		result["nonce"]=strconv.Itoa(tx.Nonce)
	default: 
		result["error"]="Command not found"
	}	
	res.Callno=request.Callno
	res.Com=request.Com
	res.Rmsg=request.Rmsg
	res.Result=result
	resonly.Result=res.Result
	if resBool {
		fmt.Println(resonly)
	} else {
		fmt.Println(res)
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
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
			if resBool {
				WriteJSON(w,resonly)
			} else {
				WriteJSON(w,result)
			}
        default:
			w.WriteHeader(http.StatusMethodNotAllowed)
 			fmt.Println("Invalid access")
        }
}

func ClientRun() {
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
		jsonRes=lib.ByteToStr(b)
	} else {
		b,err=json.Marshal(res)
		jsonRes=lib.ByteToStr(b)
	}
	fmt.Println(jsonRes)
}

// PORT NUMBER Config 추가
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
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
	return
}


// PORT NUMBER Config 연결
func HttpServer() {
	http.HandleFunc("/", apiHandler)
	log.Printf("Cubechain Http RPC Start.:"+strconv.Itoa(Configure.Httpport))
	http.ListenAndServe(":"+strconv.Itoa(Configure.Httpport), nil)
}

func ServerRun() {
	go HttpServer()
	RpcServer()
}