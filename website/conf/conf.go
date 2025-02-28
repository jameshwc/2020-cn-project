package conf

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Database struct {
	User        string
	Password    string
	Host        string
	Type        string
	Name        string
	TablePrefix string
}

type Server struct {
	RunMode   string
	HttpPort  int
	JwtSecret string
	Revision  string
}

type Redis struct {
	Host        string
	Port        int
	Password    string
	MinIdleConn int
	IdleTimeout time.Duration
}

type Log struct {
	IsLogStashActivate bool
	LogStashAddr       string
}

var DBconfig = &Database{}
var ServerConfig = &Server{}
var RedisConfig = &Redis{}
var LogConfig = &Log{}

func Setup() {
	DBconfig = &Database{
		Host: os.Getenv("db_host"),
		Type: os.Getenv("db_type"),
		Name: os.Getenv("db_name"),
	}
	port, err := strconv.Atoi(os.Getenv("server_port"))
	if err != nil {
		log.Fatal("read config: server_port is not a number: ", err.Error())
	}
	ServerConfig = &Server{
		HttpPort: port,
	}
	isLogstashActivate := false
	if os.Getenv("is_logstash_activate") == "1" {
		isLogstashActivate = true
	}
	LogConfig = &Log{
		LogStashAddr:       os.Getenv("logstash_addr"),
		IsLogStashActivate: isLogstashActivate,
	}
}
