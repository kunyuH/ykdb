package registry

import (
	"../server"
	"encoding/json"
	"fmt"
	"net/rpc"
	"os"
	"sync"
	"time"
)

var m sync.RWMutex

func Registry() {
	Data := DataInit()
	Expires := make(map[string]int64)

	DataServer := new(server.DataServer)
	DataServer.Datas = &Data
	DataServer.Expires = &Expires
	DataServer.M = &m
	rpc.Register(DataServer)

	//定时处理过期的数据
	go func(m *sync.RWMutex, Data *map[string]string,Expires *map[string]int64) {
		ti := time.NewTimer(time.Second * 2)
		for{
			<- ti.C
			//time.Sleep(2*time.Second)
			(*m).RLock()
			for k, v := range *Expires {
				if (time.Now().UnixNano() / 1e6) >= v && v != -1 && v != 0{
					(*m).Lock()
					delete(*Data,k)
					delete(*Expires,k)
					(*m).Unlock()
				}
			}
			(*m).RUnlock()
			ti.Reset(time.Second * 2)
		}
	}(&m,&Data, &Expires)
}

func DataInit() (siteinfos map[string]string) {
	//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := "C:\\MyProject\\go\\ykdb"
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

