package config

import (
    "encoding/json"
    "os"
    "fmt"
)

type Configuration struct {
	/* Network/Node */
	Network string `json:"network"`
	Nettype string `json:"nettype"`
	Chaintype string `json:"chaintype"`
	Host string `json:"host"`
	Port int `json:"port"`

	// RPC Port
	Rpcip string `json:"rpcip"`
	Rpcport int `json:"rpcport"`
	Httpport int `json:"httpport"`

	/* Cube Setting */
	Datafolder	string `json:"datafolder"`
	Datanumber	int `json:"datanumber"`

	Blocktime int `json:"blocktime"`
	Number []int `json:"number"`
	Pow []int `json:"pow"`
	Maxnonce int `json:"maxnonce"`
	Zeronumber int `json:"zeronumber"`

	Indexing int `json:"indexing"`
	Statistics int `json:"statistics"`
	Escrow int `json:"escrow"`
	Format int `json:"format"`
	Edit int `json:"edit"`

	/* Wallet */
	Address string `json:"address"`
	Keylen  int `json:"keylen"`

	/* Password */
	Password int `json:"password"`

	/* Mining */

	MiningMode string `json:"miningmode"`
	MainServer string `json:"mainserver"`
	PoolServer string `json:"poolserver"`
	PosServer string `json:"posserver"`
}

func LoadConfiguration(File string) Configuration {
    var Config Configuration
    configFile, err:=os.Open(File)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser:=json.NewDecoder(configFile)
    jsonParser.Decode(&Config)

	if Config.Network=="mainnet" {
		Config.Datanumber=10000
		Config.Maxnonce=100000000
		Config.Zeronumber=5
		Config.Blocktime=30
 		Config.Indexing=24
 		Config.Statistics=25
 		Config.Escrow=26
 		Config.Format=27
		Config.Edit=27
  		Config.Keylen=34
	} else {
		if Config.vaildConfiguration()==false {
	        panic("[Configuration Error] Please confirm configuration file.")
		}
	}
    return Config
}

func (c Configuration) vaildConfiguration() bool {
	if c.Network=="" || c.Nettype=="" || c.Host=="" || c.Port<0 {
        fmt.Println("[Configuration Error] Please confirm network infomation in configuration file.")
		return false
	}
	if c.Blocktime<10 {
        fmt.Println("[Configuration Error] Please confirm blocktime in configuration file. (Blocktime must over 10.)")
		return false
	}
	if (c.Indexing<0 || c.Statistics<0 || c.Escrow<0 || c.Format<0  || c.Edit<0) || (c.Indexing>27 || c.Statistics>27 || c.Escrow>27 || c.Format>27  || c.Edit>27) {
        fmt.Println("[Configuration Error] Please confirm special block number in configuration file. (0~27)")
		return false
	}
	return true
}
