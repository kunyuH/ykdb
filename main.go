package main

import (
	"./common"
	"./registry"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

var c *common.ConfigServer

func main(){
	//加载配置
	c = common.NewConfig()

	//注册服务
	registry.Registry()

	sock, err := net.Listen(c.Vs("net","net_work"), ":"+c.Vs("net","address"))
	log.Println("listen at :"+c.Vs("net","address"))
	if err != nil {
		log.Fatal("listen error:", err)
	}

	//c := make(chan os.Signal, 0)
	//signal.Notify(c)
	//// Block until a signal is received.
	//s := <-c
	//fmt.Println("Got signal:", s)

	for {
		conn, err := sock.Accept()
		if err != nil {
			continue
		}
		go jsonrpc.ServeConn(conn)
	}


}