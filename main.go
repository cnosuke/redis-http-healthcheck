package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/k0kubun/pp"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (r *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)

}

type ServerConfig struct {
	Bind string `yaml:"bind"`
	Port string `yaml:"port"`
}

func (s *ServerConfig) GetBinding() string {
	return fmt.Sprintf("%s:%s", s.Bind, s.Port)
}

type Config struct {
	RedisConfig  *RedisConfig  `yaml:"redis"`
	ServerConfig *ServerConfig `yaml:"server"`
}

var (
	defaultServerBind = "127.0.0.1"
	defaultServerPort = "80"
	defaultRedisHost  = "127.0.0.1"
	defaultRedisPort  = "6379"
)

func loadConfig(path string) (c *Config) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(c.RedisConfig.Host) == 0 {
		c.RedisConfig.Host = defaultRedisHost
	}

	if len(c.RedisConfig.Port) == 0 {
		c.RedisConfig.Port = defaultRedisPort
	}

	if len(c.ServerConfig.Bind) == 0 {
		c.ServerConfig.Bind = defaultServerBind
	}

	if len(c.ServerConfig.Port) == 0 {
		c.ServerConfig.Port = defaultServerPort
	}

	return c
}

func newRedisClient(c *RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.GetAddr(),
		Password: c.Password,
		DB:       c.DB,
	})
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	pong, err := redisClient.Ping().Result()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "ServiceUnavailable: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, pong)
	return
}

var (
	Name     string
	Version  string
	Revision string

	configPath  string
	config      *Config
	redisClient *redis.Client
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s (%s)", Version, Revision)
	app.Name = Name
	app.Usage = "HTTP health-check endpoint of Redis"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Path to config YAML file",
			Value:       "",
			Destination: &configPath,
		},
	}

	app.Action = func(c *cli.Context) error {
		if len(configPath) == 0 {
			return cli.NewExitError("--config should be set.", 1)
		}

		config = loadConfig(configPath)
		redisClient = newRedisClient(config.RedisConfig)

		http.HandleFunc("/healthz", httpHandler)

		pp.Printf("Starting healthcheck endpoint server...\n")
		pp.Printf("Config: %s\n", configPath)
		pp.Println(config)

		err := http.ListenAndServe(config.ServerConfig.GetBinding(), nil)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		return nil
	}

	app.Run(os.Args)
}
