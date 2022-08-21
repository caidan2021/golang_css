/*
 * @Date: 2022-08-18 11:15:22
 */
package util

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

var _config string

var Config *cnf

type cnf struct {
	MysqlHOST       string `json:"mysql_host"`
	MysqlUser       string `json:"mysql_user"`
	MysqlPassword   string `json:"mysql_password"`
	MysqlPort       string `json:"mysql_port"`
	MysqlDatabase   string `json:"mysql_database"`
	RedisHost       string `json:"redis_host"`
	RedisPort       string `json:"redis_port"`
	RedisAuth       string `json:"redis_auth"`
	RedisDB         int    `json:"redis_db"`
	HttpAddr        string `json:"http_addr"`
	FileStoragePath string `json:"file_storage_path"`
	Assets          string `json:"assets"`

	Debug bool   `json:"debug"`
	Env   string `json:"env"`
}

const (
	HttpService = "http"
)

func init() {
	flag.StringVar(&_config, "c", "./config.json", "config path")
}
func InitConfig() {
	// 初始化配置文件
	var c = &cnf{}
	f, err := os.Open(_config)
	if err != nil {
		Log.Fatal("failed to load config file, please check : is json file ? , reason : [%v]", err)
	}
	defer f.Close()

	j, err := ioutil.ReadAll(f)
	if err != nil {
		Log.Fatal("failed to read config content : %v ", err)
	}

	err = json.Unmarshal([]byte(j), c)
	if err != nil {
		Log.Fatal("failed to parse config content : %v ", err)
	}
	Config = c
}
