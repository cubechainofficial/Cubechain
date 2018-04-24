package rpc

import (
        "encoding/json"
        "fmt"
        "net/http"
        "../lib"
)

type Request struct {
	Callno	int			`json:"callno"`
	Com		string		`json:"com"`
	Vars	[]string	`json:"vars"`
	Rmsg	string		`json:"rmsg"`
}

type Response struct {
	Callno	int			`json:"callno"`
	Com		string		`json:"com"`
	Rmsg	string		`json:"rmsg"`
	Result	[]string	`json:"result"`
}

type NodeConf struct {
	Network	string `json:"network"`
	NodeID	string `json:"node"`
	Status	string `json:"status"`	
}

type PeerConf struct {
	PeerCount	string `json:"peer_count"`
	MainIp		string `json:"main_network_ip"`
	Sync		string `json:"sync"`	
	CHeight		string `json:"current_height"`	
}

type PowConf struct {
	Pows		string `json:"Pow"`
	Hashrate	string `json:"hashrate"`
	Address		string `json:"address"`	
}

type BaseConf struct {
	Rpcver		string `json:"rpcver"`
	Posstate	string `json:"posstate"`
	Cheight		string `json:"cheight"`	
}

var getError=Response{0,"Error","",[]string{"Please use post method"}}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(b)
	return nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
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
			result:=&Response{p.Callno,p.Com,p.Rmsg,p.Vars}
			ronly:=false
			result.Result,ronly=apiCommand(p.Com,p.Vars)
			if ronly {
				fmt.Println(r.Body)
				fmt.Println(result)
				WriteJSON(w,result)
			} else {
				fmt.Println(r.Body)
				fmt.Println(result)
				WriteJSON(w,result)
			}
        default:
			w.WriteHeader(http.StatusMethodNotAllowed)
 			fmt.Println("Invalid access")
        }
}

func apiCommand(command string,vars []string) ([]string,bool) {
	var Node NodeConf
	var Peer PeerConf
	var Pow PowConf
	var Base BaseConf

	var result []string
	result_only:=false
	for _,v:= range vars {
		if v=="result_only" {
			result_only=true
		}
	}
	switch command {
	case "rpc_ver":
		result=append(result,Base.Rpcver)
	case "network_info":
		result=append(result,"network:"+Node.Network)
		result=append(result,"node:"+Node.NodeID)
		result=append(result,"status:"+Node.Status)
	case "p2p_info":
		result=append(result,"peer_count:"+Peer.PeerCount)
		result=append(result,"main_network_ip:"+Peer.MainIp)
		result=append(result,"sync:"+Peer.Sync)
		result=append(result,"current_height:"+Peer.CHeight)
	case "cube_pow":
		result=append(result,"Pow:"+Pow.Pows)
		result=append(result,"Hashrate:"+Pow.Hashrate)
		result=append(result,"Address:"+Pow.Address)
	case "cube_pos":
		result=append(result,Base.Posstate)
	case "cube_height":
		result=append(result,Base.Cheight)
	case "cube_balance":
		result=append(result,lib.CallRpc(command))
	case "cube_transaction_count":
		result=append(result,lib.CallRpc(command))
	case "cube_transaction_list":
		result=append(result,lib.CallRpc(command))
	case "cube_transaction_detail":
		result=append(result,lib.CallRpc(command))
	case "cube_transaction_data":
		result=append(result,lib.CallRpc(command))
	default: 
		result=append(result,"Command not found")
	}
	return result,result_only
}

func RunRpc() {
	http.HandleFunc("/", apiHandler)
	fmt.Println("Cubechain API Start!")
	http.ListenAndServe(":8080", nil)
}
