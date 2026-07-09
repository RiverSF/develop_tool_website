package config

import (
	"os"
	"strconv"
	"strings"

	"develop_tools/pkg/path"
	"gopkg.in/ini.v1"
)

type appConfig struct {
	RunMode   string
	Host      string
	LocalHost string
}

type serverConfig struct {
	ServerPort   int
	ReadTimeout  int
	WriteTimeout int
}

type mysqlConfig struct {
	MysqlHost     string
	MysqlPort     int
	MysqlUser     string
	MysqlPassword string
	MysqlDb       string
}

var (
	AppConfig    *appConfig
	ServerConfig *serverConfig
	MysqlConfig  *mysqlConfig
)

func Init() error {
	cfg, err := ini.Load(configPath())
	if err != nil {
		return err
	}

	AppConfig = &appConfig{
		RunMode: envOrDefault("RUN_MODE", cfg.Section("").Key("RUN_MODE").String()),
	}

	ServerConfig = &serverConfig{
		ServerPort:   envIntOrDefault("HTTP_PORT", cfg.Section("server").Key("HTTP_PORT").MustInt(9080)),
		ReadTimeout:  envIntOrDefault("READ_TIMEOUT", cfg.Section("server").Key("READ_TIMEOUT").MustInt(60)),
		WriteTimeout: envIntOrDefault("WRITE_TIMEOUT", cfg.Section("server").Key("WRITE_TIMEOUT").MustInt(60)),
	}

	MysqlConfig = &mysqlConfig{
		MysqlHost:     envOrDefault("MYSQL_HOST", cfg.Section("mysql").Key("MYSQL_HOST").String()),
		MysqlPort:     envIntOrDefault("MYSQL_PORT", cfg.Section("mysql").Key("MYSQL_PORT").MustInt(3306)),
		MysqlUser:     envOrDefault("MYSQL_USER", cfg.Section("mysql").Key("MYSQL_USER").String()),
		MysqlPassword: envOrDefault("MYSQL_PASSWORD", cfg.Section("mysql").Key("MYSQL_PASSWORD").String()),
		MysqlDb:       envOrDefault("MYSQL_DB", cfg.Section("mysql").Key("MYSQL_DB").String()),
	}

	AppConfig.Host = envOrDefault("APP_HOST", cfg.Section("host").Key("HOST_"+strings.ToUpper(AppConfig.RunMode)).String())
	if AppConfig.Host == "" {
		AppConfig.Host = "http://127.0.0.1:" + strconv.Itoa(ServerConfig.ServerPort)
	}
	AppConfig.LocalHost = AppConfig.Host

	return nil
}

func configPath() string {
	if p := os.Getenv("CONFIG_PATH"); p != "" {
		return p
	}
	return path.Join("internal", "conf", "app.ini")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
