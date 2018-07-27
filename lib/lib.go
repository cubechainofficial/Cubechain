package lib

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"time"
    "os"
)



func Err(err error, exit int) int {
	if err != nil {
		fmt.Println(err)	
	}
	if exit>=1 {
		os.Exit(exit)
		return 1
	}
	return 0
}


func StrToByte(str string) []byte {
	sb := make([]byte, len(str))
	for k, v := range str {
		sb[k] = byte(v)
	}
	return sb[:]
}

func ByteToStr(bytes []byte) string {
	var str []byte
	for _, v := range bytes {
		if v != 0x0 {
			str = append(str, v)
		}
	}
	return fmt.Sprintf("%s", str)
}

func ArrayOfBytes(i int, b byte) (p []byte) {

	for i != 0 {

		p = append(p, b)
		i--
	}
	return
}

func FitBytesInto(d []byte, i int) []byte {

	if len(d) < i {

		dif := i - len(d)

		return append(ArrayOfBytes(dif, 0), d...)
	}

	return d[:i]
}

func StripByte(d []byte, b byte) []byte {

	for i, bb := range d {

		if bb != b {
			return d[i:]
		}
	}

	return nil
}

func IsNil(v interface{}) bool {
	return reflect.ValueOf(v).IsNil()
}

func DecodeJSON(r io.Reader, t interface{}) (err error) {
	err = json.NewDecoder(r).Decode(t)
	return
}

func SHA1(data []byte) string {
	hash := sha1.New()
	hash.Write(data)
	return SHAString(hash.Sum(nil))
}

func SHAString(data []byte) string {
	return fmt.Sprintf("%x", data)
}


func Timeout(i time.Duration) chan bool {
	t := make(chan bool)
	go func() {
		time.Sleep(i)
		t <- true
	}()
	return t
}


func CallRpc(com string,vars []string) string {
	for k,v:= range vars {
		fmt.Println(k,v)		
	}
	switch com {
		case "cube_balance":
		case "cube_transaction_count":
		case "cube_transaction_list":
		case "cube_transaction_detail":
		case "cube_transaction_data":
	}
	
	return fmt.Sprintf("%x", com)
}



