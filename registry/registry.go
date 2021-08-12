package registry

import (
	"../server"
	"../common"
	"encoding/json"
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
	"time"
)


func Registry() {
	Data := DataInit()
	Expires := make(map[string]int64)

	DataServer := new(server.DataServer)
	DataServer.Datas = &Data
	DataServer.Expires = &Expires
	rpc.Register(DataServer)

	//计时器
	go func() {
		var m sync.RWMutex
		for{
			time.Sleep(1*time.Second)
			for k, v := range Expires {
				if time.Now().UnixNano() / 1e6 >= v{
					m.Lock()
					delete(Data,k)
					delete(Expires,k)
					m.Unlock()
				}
			}
		}
	}()
}

func DataInit() (siteinfos map[string]string) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//dir := "C:\\MyProject\\go\\test\\kredis"
	ptjsonpath:=dir + "/data/data.json"

	f, err := os.Open(ptjsonpath)
	if err != nil {
		fmt.Println("open file err = ", err)
		return
	}

	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&siteinfos)
	if err != nil {
		fmt.Printf("json decode has error:%v\n", err)
	}
	return
}

