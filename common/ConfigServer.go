package common

import (
	"github.com/go-ini/ini"
)

type ConfigServer struct {
	cfg *ini.File
}

func (this *ConfigServer) Vs(region string,key string) string {
	return this.cfg.Section(region).Key(key).String()
}
func (this *ConfigServer) Vi(region string,key string) int {
	value,_ := this.cfg.Section(region).Key(key).Int()
	return value;
}

//实例化
func NewConfig() *ConfigServer {
	//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir := "C:\\MyProject\\go\\ykdb"
	cfg, _ := ini.Load(dir + "/config/kredis.ini")

	return &ConfigServer{
		cfg: cfg,
	}
}
