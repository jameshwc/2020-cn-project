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
	LogStashAddr string
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
		RunMode:   os.Getenv("server_runmode"),
		HttpPort:  port,
		JwtSecret: os.Getenv("server_jwt_secret"),
		Revision:  os.Getenv("server_git_commit_sha"),
	}
	minIdleConn, err := strconv.Atoi(os.Getenv("redis_min_idle_conn"))
	if err != nil {
		log.Fatal("read config: redis_min_idle_conn is not a number")
	}
	idleTimeout, err := strconv.Atoi(os.Getenv("redis_idle_timeout"))
	if err != nil {
		log.Fatal("read config: redis_idle_timeout is not a number")
	}
	redisPort, err := strconv.Atoi(os.Getenv("redis_port"))
	if err != nil {
		log.Fatal("read config: redis_port is not a number")
	}
	RedisConfig = &Redis{
		Host:        os.Getenv("redis_host"),
		Port:        redisPort,
		Password:    os.Getenv("redis_password"),
		MinIdleConn: minIdleConn,
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
	}
	LogConfig = &Log{
		LogStashAddr: os.Getenv("logstash_addr"),
	}
}
