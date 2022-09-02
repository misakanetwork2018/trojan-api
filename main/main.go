package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"trojan-api/action"
	"trojan-api/utils"
)

var (
	accessKey string
)

type Config struct {
	Web    WebApiConfig       `yaml:"web"`
	Trojan utils.TrojanConfig `yaml:"trojan"`
}

type WebApiConfig struct {
	AccessKey string `yaml:"access-key"`
	Address   string `yaml:"address"`
	Port      string `yaml:"port"`
}

func main() {
	var configFile string

	flag.StringVar(&configFile, "c", "/etc/trojan/api.yaml", "Config file path, default: /etc/trojan/api.yaml")
	flag.Parse()

	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var conf Config

	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if conf.Trojan.BinFile == "" || conf.Trojan.ApiAddr == "" || conf.Trojan.ApiPort == "" {
		fmt.Println("You should provide Trojan Api Config")
		os.Exit(1)
	}

	if conf.Web.AccessKey == "" {
		fmt.Println("You should provide a key config")
		os.Exit(1)
	}

	var address = conf.Web.Address
	if address == "" {
		address = "127.0.0.1"
	}

	var port = conf.Web.Port
	if port == "" {
		port = "8080"
	}

	utils.TrojanConf = conf.Trojan
	accessKey = conf.Web.AccessKey

	r := gin.Default()
	r.Use(webMiddleware)

	r.POST("/add", action.AddUser())
	r.POST("/del", action.DelUser())
	r.GET("/list", action.List())
	r.GET("/user", action.GetUser())
	r.POST("/user", action.ModifyUser())
	r.GET("/status", action.Status())

	_ = r.Run(address + ":" + port)
}

func webMiddleware(c *gin.Context) {
	token := c.GetHeader("X-Auth-Token")
	if token == "" {
		utils.RespondWithError(401, "API token required", c)
		return
	}
	if token != accessKey {
		utils.RespondWithError(401, "API token incorrect", c)
		return
	}
	c.Next()
}
