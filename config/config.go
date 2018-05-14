package config

import (
    "encoding/json"
    "os"
    "fmt"
)

type Configuration struct {
	Network string `json:"network"`
    Nettype string `json:"nettype"`
    Chaintype string `json:"chaintype"`
	Host string `json:"host"`
    Port int `json:"port"`
	Blocktime int `json:"blocktime"`
    Number []int `json:"number"`
    Pow []int `json:"pow"`
    Indexing int `json:"indexing"`
    Statistics int `json:"statistics"`
    Escrow int `json:"escrow"`
    Format int `json:"format"`
    Edit int `json:"edit"`
	Keylen  int `json:"keylen"`
	Password []int `json:"password"`
}

func LoadConfiguration(File string) Configuration {
    var Config Configuration
    configFile, err := os.Open(File)
    defer configFile.Close()
    if err != nil {
        fmt.Println(err.Error())
    }
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&Config)

	if Config.Network=="mainnet" {
		Config.Blocktime=180
 		Config.Number=[]int{1,2000,10000,30000,40000,50000}
 		Config.Pow=[]int{7,6,5,4,3,2,1}
 		Config.Indexing=14
 		Config.Statistics=5
 		Config.Escrow=23
 		Config.Format=0
		Config.Edit=0
  		Config.Keylen=38
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

