package server

import (
	"fmt"
	"math"
	"sync"
	"time"
	"../common"
)

//定义一个服务
type DataServer struct {
	Datas *map[string]string
	Expires *map[string]int64		//key  过期时间（毫秒）
	NewTime int64
	Config common.ConfigServer
}

//定义服务所需参数
type KeyData struct {
	Key string
}

type SetData struct {
	Key string
	Data string
	Ex int
}
type SetEx struct {
	Key string
	Ex int
}
type None struct {}

var (
	mu  sync.RWMutex
	key string
)

/**
获取key的值
*/
func (this *DataServer) Get (args KeyData,data *string) error {
	this.showNum()
	behaviors(this,args.Key)
	*data = (*this.Datas)[args.Key]
	return nil
}

func (this *DataServer) Set(args SetData,data *int) error {
	defer this.showData()
	mu.Lock()
	//设置数据
	(*this.Datas)[args.Key] = args.Data
	//true 设置过期时间
	if args.Ex == -1{
		(*this.Expires)[args.Key] = -1
	}else{
		NewTime := time.Now().UnixNano() / 1e6
		NewTime = NewTime+int64(args.Ex*1000)
		(*this.Expires)[args.Key] = NewTime
	}
	mu.Unlock()
	*data = 1
	return nil
}

/**
获取key的过期时间
key没找到则返回-1
 */
func (this *DataServer) Ttl(args KeyData,data *int) error {
	defer this.showNum()
	behaviors(this,args.Key)

	if (*this.Datas)[args.Key] == ""{
		*data = -1
		return nil
	}

	Expires := (*this.Expires)[args.Key]
	if Expires == -1 || Expires == 0{
		*data = -1
	}else{
		*data =int(math.Floor(float64((Expires-this.NewTime)/1000)))
	}
	return nil
}
/**
设置key的过期时间
key没找到则返回-1
 */
func (this *DataServer) Explre(args SetEx,data *int) error {
	defer this.showNum()
	behaviors(this,args.Key)
	if (*this.Datas)[args.Key] == ""{
		*data = -1
	}else {
		if args.Ex == -1{
			(*this.Expires)[args.Key] = -1
		}else{
			(*this.Expires)[args.Key] = this.NewTime + int64(args.Ex*1000)
		}
		*data = 1
	}
	return nil
}

func (this *DataServer) Del(args KeyData,data *int) error {
	defer this.showData()
	mu.Lock()
	delete(*this.Datas, args.Key)
	delete(*this.Expires, args.Key)
	mu.Unlock()
	*data = 1
	this.showNum()
	return nil
}

func (this *DataServer) DelAll(_ None,data *int) error {
	defer this.showData()
	mu.Lock()
	NoneMap := make(map[string]string)
	NoneMapE := make(map[string]int64)
	*this.Datas = NoneMap
	*this.Expires = NoneMapE
	mu.Unlock()
	*data = 1
	this.showNum()
	return nil
}

func (this *DataServer) List(args SetData,data *map[string]string) error {
	*data = *this.Datas
	return nil
}

func (this *DataServer) Count(_ None,data *int) error {
	*data = len(*this.Datas)
	return nil
}
func (this *DataServer) Debug(_ None,data *int) error {
	for k, v := range *this.Datas {
		fmt.Println("newTime",time.Now().UnixNano() / 1e6)
		fmt.Println("k:",k,"==v:",v)
		fmt.Println("EX:",(*this.Expires)[k])
		fmt.Println("======")
	}
	*data = 0
	return nil
}

/**
执行前执行
 */
func behaviors(this *DataServer,key string) {
	//校验权限
	//AuthServer := &AuthServer{
	//	user:this.Config.Vs("user","user"),
	//	password:this.Config.Vs("user","password"),
	//}
	//return AuthServer.check(user,password)

	//验证key的有效性
	CheckKey(this,key)
}


func CheckKey(this *DataServer,key string)  {
	Expires := (*this.Expires)[key]
	this.NewTime = time.Now().UnixNano() / 1e6
	if Expires == 0{
		(*this.Expires)[key] = -1
	}else if Expires != -1 {
		//true 已过期
		if this.NewTime >= Expires {
			mu.Lock()
			delete(*this.Datas, key)
			delete(*this.Expires, key)
			mu.Unlock()
		}
	}
}

func (this *DataServer) showData()  {
	fmt.Println(*this.Datas)
}

func (this *DataServer) showNum()  {
	fmt.Println(len(*this.Datas))
}