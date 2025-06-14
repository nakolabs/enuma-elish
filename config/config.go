package config

import (
	"encoding/json"
	"io"
	"os"
)

type App struct {
	Name string `json:"name"`
	Env  string `json:"env"`
}

type Postgres struct {
	Host            string `json:"host"`
	Port            string `json:"port"`
	Database        string `json:"database"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Driver          string `json:"driver"`
	SSLMode         string `json:"ssl_mode"`
	MaxOpenConn     int    `json:"max_open_conn"`
	MaxIdleConn     int    `json:"max_idle_conn"`
	MaxConnLifetime int    `json:"max_conn_lifetime"`  //minute
	MaxConnIdleTime int    `json:"max_conn_idle_time"` // minute
}

type Http struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
	FrontendHost string `json:"frontend_host"`
}

type JWT struct {
	Secret   string `json:"secret"`
	Duration int    `json:"duration"`
}

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
	Username string `json:"username"`
}

type SMTP struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Telemetry struct {
	Enable           bool   `json:"enable"`
	ServiceName      string `json:"service_name"`
	ExporterEndpoint string `json:"exporter_endpoint"`
	SecureMode       bool   `json:"secure_mode"`
}

type Cloudinary struct {
	CloudName string `json:"cloud_name"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	Folder    string `json:"folder"`
}

type Config struct {
	App        App        `json:"app"`
	Http       Http       `json:"http"`
	Postgres   Postgres   `json:"postgres"`
	JWT        JWT        `json:"jwt"`
	Redis      Redis      `json:"redis"`
	SMTP       SMTP       `json:"smtp"`
	Cloudinary Cloudinary `json:"cloudinary"`
	Telemetry  Telemetry  `json:"telemetry"`
}

func New(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
